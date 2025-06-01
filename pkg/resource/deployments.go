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

func (h *Handler) registerDeployments(m *server.MCPServer) {
	m.AddResource(mcp.NewResource("k8s://deployments", "Deployments",
		mcp.WithResourceDescription("List and view deployments across all namespaces"),
		mcp.WithMIMEType("application/json"),
	), h.getDeployments)
	m.AddResourceTemplate(mcp.NewResourceTemplate(
		"k8s://{namespace}/deployments",
		"Deployments in namespace",
		mcp.WithTemplateDescription("List and view deployments in a specific namespace"),
		mcp.WithTemplateMIMEType("application/json"),
	), h.getDeploymentsInNamespace)
}

func (h *Handler) getDeployments(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	deployments, err := h.client.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments in all namespaces: %w", err)
	}

	var summaries []map[string]interface{}
	for _, d := range deployments.Items {
		summaries = append(summaries, map[string]interface{}{
			"name":      d.Name,
			"namespace": d.Namespace,
			"replicas":  d.Status.Replicas,
			"available": d.Status.AvailableReplicas,
			"age":       time.Since(d.CreationTimestamp.Time).Round(time.Second).String(),
		})
	}

	result, err := json.MarshalIndent(map[string]interface{}{
		"scope":       "All namespaces",
		"total_items": len(summaries),
		"deployments": summaries,
	}, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal deployment summaries: %w", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(result),
		},
	}, nil
}

func (h *Handler) getDeploymentsInNamespace(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	uri := request.Params.URI
	ns, _ := ExtractNamespaceFromURI(uri)
	deployments, err := h.client.AppsV1().Deployments(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments in namespace '%s': %w", ns, err)
	}

	var summaries []map[string]interface{}
	for _, d := range deployments.Items {
		summaries = append(summaries, map[string]interface{}{
			"name":      d.Name,
			"namespace": d.Namespace,
			"replicas":  d.Status.Replicas,
			"available": d.Status.AvailableReplicas,
			"age":       time.Since(d.CreationTimestamp.Time).Round(time.Second).String(),
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
		"deployments": summaries,
	}, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal deployment summaries: %w", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      uri,
			MIMEType: "application/json",
			Text:     string(result),
		},
	}, nil
}
