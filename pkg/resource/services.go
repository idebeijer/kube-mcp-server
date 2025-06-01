package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (h *Handler) registerServices(m *server.MCPServer) {
	m.AddResource(mcp.NewResource("k8s://services", "Services",
		mcp.WithResourceDescription("List and view services across all namespaces"),
		mcp.WithMIMEType("application/json"),
	), h.getServices)
	m.AddResourceTemplate(mcp.NewResourceTemplate(
		"k8s://{namespace}/services",
		"Services in namespace",
		mcp.WithTemplateDescription("List and view services in a specific namespace"),
		mcp.WithTemplateMIMEType("application/json"),
	), h.getServicesInNamespace)
}

func (h *Handler) getServices(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	services, err := h.client.CoreV1().Services("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services in all namespaces: %w", err)
	}

	var summaries []map[string]interface{}
	for _, s := range services.Items {
		ports := make([]string, 0, len(s.Spec.Ports))
		for _, p := range s.Spec.Ports {
			ports = append(ports, fmt.Sprintf("%d/%s", p.Port, p.Protocol))
		}
		summaries = append(summaries, map[string]interface{}{
			"name":      s.Name,
			"namespace": s.Namespace,
			"type":      s.Spec.Type,
			"clusterIP": s.Spec.ClusterIP,
			"ports":     ports,
			"age":       time.Since(s.CreationTimestamp.Time).Round(time.Second).String(),
		})
	}

	result, err := json.MarshalIndent(map[string]interface{}{
		"scope":       "All namespaces",
		"total_items": len(summaries),
		"services":    summaries,
	}, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal service summaries: %w", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(result),
		},
	}, nil
}

func (h *Handler) getServicesInNamespace(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	uri := request.Params.URI
	ns, _ := ExtractNamespaceFromURI(uri)
	services, err := h.client.CoreV1().Services(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services in namespace '%s': %w", ns, err)
	}

	var summaries []map[string]interface{}
	for _, s := range services.Items {
		ports := make([]string, 0, len(s.Spec.Ports))
		for _, p := range s.Spec.Ports {
			ports = append(ports, fmt.Sprintf("%d/%s", p.Port, p.Protocol))
		}
		summaries = append(summaries, map[string]interface{}{
			"name":      s.Name,
			"namespace": s.Namespace,
			"type":      s.Spec.Type,
			"clusterIP": s.Spec.ClusterIP,
			"ports":     ports,
			"age":       time.Since(s.CreationTimestamp.Time).Round(time.Second).String(),
		})
	}

	var scopeDescription string
	if ns == "" {
		scopeDescription = "All namespaces"
	} else {
		scopeDescription = fmt.Sprintf("Namespace: %s", ns)
	}

	result, err := json.MarshalIndent(map[string]interface{}{
		"scope":       scopeDescription,
		"namespace":   ns,
		"total_items": len(summaries),
		"services":    summaries,
	}, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal service summaries: %w", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      uri,
			MIMEType: "application/json",
			Text:     string(result),
		},
	}, nil
}
