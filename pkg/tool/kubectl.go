package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (h *Handler) registerKubectl(m *server.MCPServer) {
	m.AddTool(mcp.NewTool("kubectl_get",
		mcp.WithDescription("Execute kubectl get command for any Kubernetes resource type with advanced filtering options"),
		mcp.WithString("resource",
			mcp.Description("The resource type to get (e.g., pods, deployments, services, nodes, etc.)"),
			mcp.Required(),
		),
		mcp.WithString("name",
			mcp.Description("Name of the specific resource to get (optional - leave empty to get all resources of the type)"),
		),
		mcp.WithString("namespace",
			mcp.Description("Namespace to query (optional - leave empty for all namespaces or cluster-scoped resources)"),
		),
		mcp.WithString("field_selector",
			mcp.Description("Field selector to filter resources (e.g., 'status.phase=Running,metadata.namespace!=kube-system')"),
		),
		mcp.WithString("label_selector",
			mcp.Description("Label selector to filter resources (e.g., 'app=nginx,version=v1.0')"),
		),
		mcp.WithString("output",
			mcp.Description("Output format: json, yaml, table, name, wide, custom-columns, jsonpath"),
			mcp.DefaultString("json"),
		),
		mcp.WithBoolean("all_namespaces",
			mcp.Description("List resources across all namespaces (equivalent to --all-namespaces)"),
			mcp.DefaultBool(false),
		),
		mcp.WithBoolean("show_labels",
			mcp.Description("Show labels in the output (equivalent to --show-labels)"),
			mcp.DefaultBool(false),
		),
		mcp.WithString("sort_by",
			mcp.Description("Sort resources by a specific field (e.g., '.metadata.name', '.status.startTime')"),
		),
		mcp.WithString("custom_columns",
			mcp.Description("Custom columns for output when output format is 'custom-columns' (e.g., 'NAME:.metadata.name,STATUS:.status.phase')"),
		),
		mcp.WithString("jsonpath",
			mcp.Description("JSONPath expression when output format is 'jsonpath' (e.g., '{.items[*].metadata.name}')"),
		),
	), mcp.NewTypedToolHandler[KubectlGetArgs](h.kubectlGetHandler()))

	m.AddTool(mcp.NewTool("kubectl_describe",
		mcp.WithDescription("Execute kubectl describe command for detailed information about Kubernetes resources"),
		mcp.WithString("resource",
			mcp.Description("The resource type to describe (e.g., pod, deployment, service, node, etc.)"),
			mcp.Required(),
		),
		mcp.WithString("name",
			mcp.Description("Name of the specific resource to describe"),
			mcp.Required(),
		),
		mcp.WithString("namespace",
			mcp.Description("Namespace of the resource (optional for cluster-scoped resources)"),
		),
	), mcp.NewTypedToolHandler[KubectlDescribeArgs](h.kubectlDescribeHandler()))

	m.AddTool(mcp.NewTool("kubectl_logs",
		mcp.WithDescription("Execute kubectl logs command to get logs from pods"),
		mcp.WithString("pod_name",
			mcp.Description("Name of the pod to get logs from"),
			mcp.Required(),
		),
		mcp.WithString("namespace",
			mcp.Description("Namespace of the pod"),
			mcp.DefaultString("default"),
		),
		mcp.WithString("container",
			mcp.Description("Container name (optional, required for multi-container pods)"),
		),
		mcp.WithBoolean("follow",
			mcp.Description("Follow the log stream (equivalent to -f flag)"),
			mcp.DefaultBool(false),
		),
		mcp.WithBoolean("previous",
			mcp.Description("Get logs from previous container instance (equivalent to -p flag)"),
			mcp.DefaultBool(false),
		),
		mcp.WithNumber("tail",
			mcp.Description("Number of lines to show from the end of the logs"),
		),
		mcp.WithString("since",
			mcp.Description("Only return logs newer than a relative duration like 5s, 2m, or 3h"),
		),
		mcp.WithString("since_time",
			mcp.Description("Only return logs after a specific date (RFC3339)"),
		),
		mcp.WithBoolean("timestamps",
			mcp.Description("Include timestamps in the output"),
			mcp.DefaultBool(false),
		),
	), mcp.NewTypedToolHandler[KubectlLogsArgs](h.kubectlLogsHandler()))
}

type KubectlGetArgs struct {
	Resource      string `json:"resource"`
	Name          string `json:"name,omitempty"`
	Namespace     string `json:"namespace,omitempty"`
	FieldSelector string `json:"field_selector,omitempty"`
	LabelSelector string `json:"label_selector,omitempty"`
	Output        string `json:"output"`
	AllNamespaces bool   `json:"all_namespaces"`
	ShowLabels    bool   `json:"show_labels"`
	SortBy        string `json:"sort_by,omitempty"`
	CustomColumns string `json:"custom_columns,omitempty"`
	JSONPath      string `json:"jsonpath,omitempty"`
}

func (h *Handler) kubectlGetHandler() mcp.TypedToolHandlerFunc[KubectlGetArgs] {
	return func(
		ctx context.Context,
		req mcp.CallToolRequest,
		args KubectlGetArgs,
	) (*mcp.CallToolResult, error) {
		cmdArgs := []string{"get", args.Resource}

		if args.Name != "" {
			cmdArgs = append(cmdArgs, args.Name)
		}
		if args.Namespace != "" && !args.AllNamespaces {
			cmdArgs = append(cmdArgs, "-n", args.Namespace)
		}
		if args.AllNamespaces {
			cmdArgs = append(cmdArgs, "--all-namespaces")
		}
		if args.FieldSelector != "" {
			cmdArgs = append(cmdArgs, "--field-selector", args.FieldSelector)
		}
		if args.LabelSelector != "" {
			cmdArgs = append(cmdArgs, "-l", args.LabelSelector)
		}

		if args.Output != "" {
			switch args.Output {
			case "custom-columns":
				if args.CustomColumns == "" {
					return mcp.NewToolResultError("custom_columns must be specified when output format is 'custom-columns'"), nil
				}
				cmdArgs = append(cmdArgs, "-o", fmt.Sprintf("custom-columns=%s", args.CustomColumns))
			case "jsonpath":
				if args.JSONPath == "" {
					return mcp.NewToolResultError("jsonpath must be specified when output format is 'jsonpath'"), nil
				}
				cmdArgs = append(cmdArgs, "-o", fmt.Sprintf("jsonpath=%s", args.JSONPath))
			default:
				cmdArgs = append(cmdArgs, "-o", args.Output)
			}
		}

		if args.ShowLabels {
			cmdArgs = append(cmdArgs, "--show-labels")
		}
		if args.SortBy != "" {
			cmdArgs = append(cmdArgs, "--sort-by", args.SortBy)
		}

		output, err := h.runKubectl(ctx, cmdArgs...)
		if err != nil {
			errorMsg := fmt.Sprintf("kubectl command failed: %v\nCommand: kubectl %s\nOutput: %s",
				err, strings.Join(cmdArgs, " "), string(output))
			return mcp.NewToolResultError(errorMsg), nil
		}

		response := string(output)
		if args.Output == "json" && len(output) > 0 {
			var jsonData interface{}
			if err := json.Unmarshal(output, &jsonData); err == nil {
				if formatted, err := json.MarshalIndent(jsonData, "", "  "); err == nil {
					response = string(formatted)
				}
			}
		}

		fullResponse := fmt.Sprintf("Command executed: kubectl %s\n\n%s", strings.Join(cmdArgs, " "), response)

		return mcp.NewToolResultText(fullResponse), nil
	}
}

type KubectlDescribeArgs struct {
	Resource  string `json:"resource"`
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

func (h *Handler) kubectlDescribeHandler() mcp.TypedToolHandlerFunc[KubectlDescribeArgs] {
	return func(
		ctx context.Context,
		req mcp.CallToolRequest,
		args KubectlDescribeArgs,
	) (*mcp.CallToolResult, error) {
		cmdArgs := []string{"describe", args.Resource, args.Name}
		if args.Namespace != "" {
			cmdArgs = append(cmdArgs, "-n", args.Namespace)
		}

		output, err := h.runKubectl(ctx, cmdArgs...)
		if err != nil {
			errorMsg := fmt.Sprintf("kubectl describe failed: %v\nCommand: kubectl %s\nOutput: %s",
				err, strings.Join(cmdArgs, " "), string(output))
			return mcp.NewToolResultError(errorMsg), nil
		}

		fullResponse := fmt.Sprintf("Command executed: kubectl %s\n\n%s", strings.Join(cmdArgs, " "), string(output))
		return mcp.NewToolResultText(fullResponse), nil
	}
}

type KubectlLogsArgs struct {
	PodName    string `json:"pod_name"`
	Namespace  string `json:"namespace"`
	Container  string `json:"container,omitempty"`
	Follow     bool   `json:"follow"`
	Previous   bool   `json:"previous"`
	Tail       int    `json:"tail,omitempty"`
	Since      string `json:"since,omitempty"`
	SinceTime  string `json:"since_time,omitempty"`
	Timestamps bool   `json:"timestamps"`
}

func (h *Handler) kubectlLogsHandler() mcp.TypedToolHandlerFunc[KubectlLogsArgs] {
	return func(
		ctx context.Context,
		req mcp.CallToolRequest,
		args KubectlLogsArgs,
	) (*mcp.CallToolResult, error) {
		cmdArgs := []string{"logs", args.PodName}

		if args.Namespace != "" {
			cmdArgs = append(cmdArgs, "-n", args.Namespace)
		}
		if args.Container != "" {
			cmdArgs = append(cmdArgs, "-c", args.Container)
		}
		if args.Follow {
			cmdArgs = append(cmdArgs, "-f")
		}
		if args.Previous {
			cmdArgs = append(cmdArgs, "-p")
		}
		if args.Tail > 0 {
			cmdArgs = append(cmdArgs, "--tail", fmt.Sprintf("%d", args.Tail))
		}
		if args.Since != "" {
			cmdArgs = append(cmdArgs, "--since", args.Since)
		}
		if args.SinceTime != "" {
			cmdArgs = append(cmdArgs, "--since-time", args.SinceTime)
		}
		if args.Timestamps {
			cmdArgs = append(cmdArgs, "--timestamps")
		}
		output, err := h.runKubectl(ctx, cmdArgs...)
		if err != nil {
			errorMsg := fmt.Sprintf("kubectl logs failed: %v\nCommand: kubectl %s\nOutput: %s",
				err, strings.Join(cmdArgs, " "), string(output))
			return mcp.NewToolResultError(errorMsg), nil
		}

		fullResponse := fmt.Sprintf("Command executed: kubectl %s\n\n%s", strings.Join(cmdArgs, " "), string(output))
		return mcp.NewToolResultText(fullResponse), nil
	}
}
