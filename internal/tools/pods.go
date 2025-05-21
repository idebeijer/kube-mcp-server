package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CountPodsToolDefinition() mcp.Tool {
	return mcp.NewTool("count_pods",
		mcp.WithDescription("Count Pods in a Kubernetes namespace"),
		mcp.WithString("namespace",
			mcp.Description("Namespace to count (empty for all)"),
			mcp.DefaultString("default"),
		),
	)
}

type CountPodsArgs struct {
	Namespace string `json:"namespace"`
}

func CountPodsHandler(client *kubernetes.Clientset) mcp.TypedToolHandlerFunc[CountPodsArgs] {
	return func(
		ctx context.Context,
		req mcp.CallToolRequest,
		args CountPodsArgs,
	) (*mcp.CallToolResult, error) {
		podsCount, err := client.CoreV1().
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
