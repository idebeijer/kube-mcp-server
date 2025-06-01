package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/idebeijer/kube-mcp-server/pkg/kube"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (h *Handler) registerPods(m *server.MCPServer) {
	m.AddResource(mcp.Resource{
		URI:         "k8s://pods",
		Name:        "Pods",
		Description: "List and view pods across all namespaces",
		MIMEType:    "application/json",
	}, h.getPods)
	m.AddResourceTemplate(mcp.NewResourceTemplate("k8s://{namespace}/pods", "Pods in namespace"), h.getPodsInNamespace)
}

func (h *Handler) getPodsInNamespace(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	uri := request.Params.URI
	targetNamespace, _ := ExtractNamespaceFromURI(uri)

	pods, err := h.client.CoreV1().Pods(targetNamespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods in namespace '%s': %w", targetNamespace, err)
	}

	var podSummaries []map[string]interface {
	}
	for _, pod := range pods.Items {
		ready, total := kube.GetPodReadyContainers(pod.Status.ContainerStatuses)
		podSummaries = append(podSummaries, map[string]interface{}{
			"name":      pod.Name,
			"namespace": pod.Namespace,
			"status":    string(pod.Status.Phase),
			"ready":     fmt.Sprintf("%d/%d", ready, total),
			"node":      pod.Spec.NodeName,
			"age":       time.Since(pod.CreationTimestamp.Time).Round(time.Second).String(),
			"containers": func() []string {
				var containers []string
				for _, container := range pod.Spec.Containers {
					containers = append(containers, container.Name)
				}
				return containers
			}(),
		})
	}

	var scopeDescription string
	if targetNamespace == "" {
		scopeDescription = "All namespaces"
	} else {
		scopeDescription = fmt.Sprintf("Namespace: %s", targetNamespace)
	}

	result, err := json.MarshalIndent(map[string]interface{}{
		"scope":      scopeDescription,
		"namespace":  targetNamespace,
		"total_pods": len(podSummaries),
		"pods":       podSummaries,
	}, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal pod summaries: %w", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      uri,
			MIMEType: "application/json",
			Text:     string(result),
		},
	}, nil
}

func (h *Handler) getPods(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	pods, err := h.client.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods in all namespace: %w", err)
	}

	var podSummaries []map[string]interface {
	}
	for _, pod := range pods.Items {
		ready, total := kube.GetPodReadyContainers(pod.Status.ContainerStatuses)
		podSummaries = append(podSummaries, map[string]interface{}{
			"name":      pod.Name,
			"namespace": pod.Namespace,
			"status":    string(pod.Status.Phase),
			"ready":     fmt.Sprintf("%d/%d", ready, total),
			"node":      pod.Spec.NodeName,
			"age":       time.Since(pod.CreationTimestamp.Time).Round(time.Second).String(),
			"containers": func() []string {
				var containers []string
				for _, container := range pod.Spec.Containers {
					containers = append(containers, container.Name)
				}
				return containers
			}(),
		})
	}

	result, err := json.MarshalIndent(map[string]interface{}{
		"scope":      "All namespaces",
		"total_pods": len(podSummaries),
		"pods":       podSummaries,
	}, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal pod summaries: %w", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(result),
		},
	}, nil
}
