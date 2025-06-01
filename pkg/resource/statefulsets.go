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

func (h *Handler) registerStatefulSets(m *server.MCPServer) {
	m.AddResource(mcp.NewResource("k8s://statefulsets", "StatefulSets",
		mcp.WithResourceDescription("List and view statefulsets across all namespaces"),
		mcp.WithMIMEType("application/json"),
	), h.getStatefulSets)
	m.AddResourceTemplate(mcp.NewResourceTemplate(
		"k8s://{namespace}/statefulsets",
		"StatefulSets in namespace",
		mcp.WithTemplateDescription("List and view statefulsets in a specific namespace"),
		mcp.WithTemplateMIMEType("application/json"),
	), h.getStatefulSetsInNamespace)
}

func (h *Handler) getStatefulSets(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	sets, err := h.client.AppsV1().StatefulSets("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list statefulsets in all namespaces: %w", err)
	}

	var summaries []map[string]interface{}
	for _, s := range sets.Items {
		summaries = append(summaries, map[string]interface{}{
			"name":      s.Name,
			"namespace": s.Namespace,
			"replicas":  s.Status.Replicas,
			"ready":     s.Status.ReadyReplicas,
			"age":       time.Since(s.CreationTimestamp.Time).Round(time.Second).String(),
		})
	}

	result, err := json.MarshalIndent(map[string]interface{}{
		"scope":        "All namespaces",
		"total_items":  len(summaries),
		"statefulsets": summaries,
	}, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal statefulset summaries: %w", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(result),
		},
	}, nil
}

func (h *Handler) getStatefulSetsInNamespace(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	uri := request.Params.URI
	ns, _ := ExtractNamespaceFromURI(uri)
	sets, err := h.client.AppsV1().StatefulSets(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list statefulsets in namespace '%s': %w", ns, err)
	}

	var summaries []map[string]interface{}
	for _, s := range sets.Items {
		summaries = append(summaries, map[string]interface{}{
			"name":      s.Name,
			"namespace": s.Namespace,
			"replicas":  s.Status.Replicas,
			"ready":     s.Status.ReadyReplicas,
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
		"scope":        scopeDescription,
		"namespace":    ns,
		"total_items":  len(summaries),
		"statefulsets": summaries,
	}, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal statefulset summaries: %w", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      uri,
			MIMEType: "application/json",
			Text:     string(result),
		},
	}, nil
}
