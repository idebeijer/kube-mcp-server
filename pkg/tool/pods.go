package tool

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (h *Handler) registerPods(m *server.MCPServer) {
	m.AddTool(mcp.NewTool("count_pods",
		mcp.WithDescription("Count Pods in a Kubernetes namespace"),
		mcp.WithString("namespace",
			mcp.Description("Namespace to count (empty for all)"),
			mcp.DefaultString("default"),
		),
	), mcp.NewTypedToolHandler[CountPodsArgs](h.countPodsHandler()))
}

type CountPodsArgs struct {
	Namespace string `json:"namespace"`
}

func (h *Handler) countPodsHandler() mcp.TypedToolHandlerFunc[CountPodsArgs] {
	return func(
		ctx context.Context,
		req mcp.CallToolRequest,
		args CountPodsArgs,
	) (*mcp.CallToolResult, error) {
		podsCount, err := h.client.CoreV1().
			Pods(args.Namespace).
			List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("count pods failed", err), nil
		}
		return mcp.NewToolResultText(
			fmt.Sprintf("Found %d pods", len(podsCount.Items)),
		), nil
	}
}
