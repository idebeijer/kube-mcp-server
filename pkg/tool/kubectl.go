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

	m.AddTool(mcp.NewTool("kubectl_create",
		mcp.WithDescription("Execute kubectl create command to create Kubernetes resources"),
		mcp.WithString("filename",
			mcp.Description("Filename or URL of the resource to create (e.g., deployment.yaml)"),
		),
		mcp.WithString("resource",
			mcp.Description("Resource type to create (e.g., deployment, service, configmap)"),
		),
		mcp.WithString("name",
			mcp.Description("Name of the resource to create"),
		),
		mcp.WithString("namespace",
			mcp.Description("Namespace to create the resource in"),
		),
		mcp.WithString("image",
			mcp.Description("Container image (for creating deployments)"),
		),
		mcp.WithBoolean("dry_run",
			mcp.Description("Run in dry-run mode without actually creating the resource"),
			mcp.DefaultBool(false),
		),
		mcp.WithString("output",
			mcp.Description("Output format: json, yaml, name"),
		),
	), mcp.NewTypedToolHandler[KubectlCreateArgs](h.kubectlCreateHandler()))

	m.AddTool(mcp.NewTool("kubectl_delete",
		mcp.WithDescription("Execute kubectl delete command to delete Kubernetes resources"),
		mcp.WithString("resource",
			mcp.Description("Resource type to delete (e.g., pod, deployment, service)"),
			mcp.Required(),
		),
		mcp.WithString("name",
			mcp.Description("Name of the resource to delete (optional - can use label selector instead)"),
		),
		mcp.WithString("filename",
			mcp.Description("Filename or URL of the resource to delete"),
		),
		mcp.WithString("namespace",
			mcp.Description("Namespace of the resource to delete"),
		),
		mcp.WithString("label_selector",
			mcp.Description("Label selector to delete multiple resources (e.g., 'app=nginx')"),
		),
		mcp.WithBoolean("all",
			mcp.Description("Delete all resources of the specified type in the namespace"),
			mcp.DefaultBool(false),
		),
		mcp.WithBoolean("force",
			mcp.Description("Force delete the resource (equivalent to --force)"),
			mcp.DefaultBool(false),
		),
		mcp.WithNumber("grace_period",
			mcp.Description("Grace period in seconds for pod deletion"),
		),
		mcp.WithBoolean("ignore_not_found",
			mcp.Description("Treat \"resource not found\" as success"),
			mcp.DefaultBool(false),
		),
	), mcp.NewTypedToolHandler[KubectlDeleteArgs](h.kubectlDeleteHandler()))

	m.AddTool(mcp.NewTool("kubectl_apply",
		mcp.WithDescription("Execute kubectl apply command to apply configuration to resources"),
		mcp.WithString("filename",
			mcp.Description("Filename, directory, or URL of the resource to apply"),
			mcp.Required(),
		),
		mcp.WithString("namespace",
			mcp.Description("Namespace to apply the resource in"),
		),
		mcp.WithBoolean("recursive",
			mcp.Description("Process the directory used in -f, --filename recursively"),
			mcp.DefaultBool(false),
		),
		mcp.WithBoolean("dry_run",
			mcp.Description("Run in dry-run mode without actually applying changes"),
			mcp.DefaultBool(false),
		),
		mcp.WithString("output",
			mcp.Description("Output format: json, yaml, name"),
		),
		mcp.WithBoolean("force",
			mcp.Description("Force apply even if the resource already exists"),
			mcp.DefaultBool(false),
		),
		mcp.WithBoolean("validate",
			mcp.Description("Validate the resource before applying"),
			mcp.DefaultBool(true),
		),
	), mcp.NewTypedToolHandler[KubectlApplyArgs](h.kubectlApplyHandler()))

	m.AddTool(mcp.NewTool("kubectl_label",
		mcp.WithDescription("Execute kubectl label command to add, update, or remove labels on resources"),
		mcp.WithString("resource",
			mcp.Description("Resource type to label (e.g., pod, node, deployment)"),
			mcp.Required(),
		),
		mcp.WithString("name",
			mcp.Description("Name of the resource to label (optional - can use label selector instead)"),
		),
		mcp.WithString("labels",
			mcp.Description("Labels to set (e.g., 'key1=value1,key2=value2' or 'key1-' to remove)"),
			mcp.Required(),
		),
		mcp.WithString("namespace",
			mcp.Description("Namespace of the resource"),
		),
		mcp.WithString("label_selector",
			mcp.Description("Label selector to label multiple resources"),
		),
		mcp.WithBoolean("overwrite",
			mcp.Description("Overwrite existing labels"),
			mcp.DefaultBool(false),
		),
		mcp.WithBoolean("all",
			mcp.Description("Label all resources of the specified type in the namespace"),
			mcp.DefaultBool(false),
		),
	), mcp.NewTypedToolHandler[KubectlLabelArgs](h.kubectlLabelHandler()))

	m.AddTool(mcp.NewTool("kubectl_annotate",
		mcp.WithDescription("Execute kubectl annotate command to add, update, or remove annotations on resources"),
		mcp.WithString("resource",
			mcp.Description("Resource type to annotate (e.g., pod, node, deployment)"),
			mcp.Required(),
		),
		mcp.WithString("name",
			mcp.Description("Name of the resource to annotate (optional - can use label selector instead)"),
		),
		mcp.WithString("annotations",
			mcp.Description("Annotations to set (e.g., 'key1=value1,key2=value2' or 'key1-' to remove)"),
			mcp.Required(),
		),
		mcp.WithString("namespace",
			mcp.Description("Namespace of the resource"),
		),
		mcp.WithString("label_selector",
			mcp.Description("Label selector to annotate multiple resources"),
		),
		mcp.WithBoolean("overwrite",
			mcp.Description("Overwrite existing annotations"),
			mcp.DefaultBool(false),
		),
		mcp.WithBoolean("all",
			mcp.Description("Annotate all resources of the specified type in the namespace"),
			mcp.DefaultBool(false),
		),
	), mcp.NewTypedToolHandler[KubectlAnnotateArgs](h.kubectlAnnotateHandler()))
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

type KubectlCreateArgs struct {
	Filename  string `json:"filename,omitempty"`
	Resource  string `json:"resource,omitempty"`
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Image     string `json:"image,omitempty"`
	DryRun    bool   `json:"dry_run"`
	Output    string `json:"output,omitempty"`
}

func (h *Handler) kubectlCreateHandler() mcp.TypedToolHandlerFunc[KubectlCreateArgs] {
	return func(
		ctx context.Context,
		req mcp.CallToolRequest,
		args KubectlCreateArgs,
	) (*mcp.CallToolResult, error) {
		var cmdArgs []string

		if args.Filename != "" {
			cmdArgs = []string{"create", "-f", args.Filename}
		} else if args.Resource != "" {
			cmdArgs = []string{"create", args.Resource}
			if args.Name != "" {
				cmdArgs = append(cmdArgs, args.Name)
			}
			if args.Image != "" && args.Resource == "deployment" {
				cmdArgs = append(cmdArgs, "--image", args.Image)
			}
		} else {
			return mcp.NewToolResultError("Either filename or resource must be specified"), nil
		}

		if args.Namespace != "" {
			cmdArgs = append(cmdArgs, "-n", args.Namespace)
		}
		if args.DryRun {
			cmdArgs = append(cmdArgs, "--dry-run=client")
		}
		if args.Output != "" {
			cmdArgs = append(cmdArgs, "-o", args.Output)
		}

		output, err := h.runKubectl(ctx, cmdArgs...)
		if err != nil {
			errorMsg := fmt.Sprintf("kubectl create failed: %v\nCommand: kubectl %s\nOutput: %s",
				err, strings.Join(cmdArgs, " "), string(output))
			return mcp.NewToolResultError(errorMsg), nil
		}

		fullResponse := fmt.Sprintf("Command executed: kubectl %s\n\n%s", strings.Join(cmdArgs, " "), string(output))
		return mcp.NewToolResultText(fullResponse), nil
	}
}

type KubectlDeleteArgs struct {
	Resource       string `json:"resource"`
	Name           string `json:"name,omitempty"`
	Filename       string `json:"filename,omitempty"`
	Namespace      string `json:"namespace,omitempty"`
	LabelSelector  string `json:"label_selector,omitempty"`
	All            bool   `json:"all"`
	Force          bool   `json:"force"`
	GracePeriod    int    `json:"grace_period,omitempty"`
	IgnoreNotFound bool   `json:"ignore_not_found"`
}

func (h *Handler) kubectlDeleteHandler() mcp.TypedToolHandlerFunc[KubectlDeleteArgs] {
	return func(
		ctx context.Context,
		req mcp.CallToolRequest,
		args KubectlDeleteArgs,
	) (*mcp.CallToolResult, error) {
		var cmdArgs []string

		if args.Filename != "" {
			cmdArgs = []string{"delete", "-f", args.Filename}
		} else {
			cmdArgs = []string{"delete", args.Resource}
			if args.Name != "" {
				cmdArgs = append(cmdArgs, args.Name)
			}
		}

		if args.Namespace != "" {
			cmdArgs = append(cmdArgs, "-n", args.Namespace)
		}
		if args.LabelSelector != "" {
			cmdArgs = append(cmdArgs, "-l", args.LabelSelector)
		}
		if args.All {
			cmdArgs = append(cmdArgs, "--all")
		}
		if args.Force {
			cmdArgs = append(cmdArgs, "--force")
		}
		if args.GracePeriod > 0 {
			cmdArgs = append(cmdArgs, "--grace-period", fmt.Sprintf("%d", args.GracePeriod))
		}
		if args.IgnoreNotFound {
			cmdArgs = append(cmdArgs, "--ignore-not-found")
		}

		output, err := h.runKubectl(ctx, cmdArgs...)
		if err != nil {
			errorMsg := fmt.Sprintf("kubectl delete failed: %v\nCommand: kubectl %s\nOutput: %s",
				err, strings.Join(cmdArgs, " "), string(output))
			return mcp.NewToolResultError(errorMsg), nil
		}

		fullResponse := fmt.Sprintf("Command executed: kubectl %s\n\n%s", strings.Join(cmdArgs, " "), string(output))
		return mcp.NewToolResultText(fullResponse), nil
	}
}

type KubectlApplyArgs struct {
	Filename  string `json:"filename"`
	Namespace string `json:"namespace,omitempty"`
	Recursive bool   `json:"recursive"`
	DryRun    bool   `json:"dry_run"`
	Output    string `json:"output,omitempty"`
	Force     bool   `json:"force"`
	Validate  bool   `json:"validate"`
}

func (h *Handler) kubectlApplyHandler() mcp.TypedToolHandlerFunc[KubectlApplyArgs] {
	return func(
		ctx context.Context,
		req mcp.CallToolRequest,
		args KubectlApplyArgs,
	) (*mcp.CallToolResult, error) {
		cmdArgs := []string{"apply", "-f", args.Filename}

		if args.Namespace != "" {
			cmdArgs = append(cmdArgs, "-n", args.Namespace)
		}
		if args.Recursive {
			cmdArgs = append(cmdArgs, "--recursive")
		}
		if args.DryRun {
			cmdArgs = append(cmdArgs, "--dry-run=client")
		}
		if args.Output != "" {
			cmdArgs = append(cmdArgs, "-o", args.Output)
		}
		if args.Force {
			cmdArgs = append(cmdArgs, "--force")
		}
		if !args.Validate {
			cmdArgs = append(cmdArgs, "--validate=false")
		}

		output, err := h.runKubectl(ctx, cmdArgs...)
		if err != nil {
			errorMsg := fmt.Sprintf("kubectl apply failed: %v\nCommand: kubectl %s\nOutput: %s",
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

type KubectlLabelArgs struct {
	Resource      string `json:"resource"`
	Name          string `json:"name,omitempty"`
	Labels        string `json:"labels"`
	Namespace     string `json:"namespace,omitempty"`
	LabelSelector string `json:"label_selector,omitempty"`
	Overwrite     bool   `json:"overwrite"`
	All           bool   `json:"all"`
}

func (h *Handler) kubectlLabelHandler() mcp.TypedToolHandlerFunc[KubectlLabelArgs] {
	return func(
		ctx context.Context,
		req mcp.CallToolRequest,
		args KubectlLabelArgs,
	) (*mcp.CallToolResult, error) {
		cmdArgs := []string{"label", args.Resource}

		if args.Name != "" {
			cmdArgs = append(cmdArgs, args.Name)
		}

		// Add labels to the command
		cmdArgs = append(cmdArgs, args.Labels)

		if args.Namespace != "" {
			cmdArgs = append(cmdArgs, "-n", args.Namespace)
		}
		if args.LabelSelector != "" {
			cmdArgs = append(cmdArgs, "-l", args.LabelSelector)
		}
		if args.Overwrite {
			cmdArgs = append(cmdArgs, "--overwrite")
		}
		if args.All {
			cmdArgs = append(cmdArgs, "--all")
		}

		output, err := h.runKubectl(ctx, cmdArgs...)
		if err != nil {
			errorMsg := fmt.Sprintf("kubectl label failed: %v\nCommand: kubectl %s\nOutput: %s",
				err, strings.Join(cmdArgs, " "), string(output))
			return mcp.NewToolResultError(errorMsg), nil
		}

		fullResponse := fmt.Sprintf("Command executed: kubectl %s\n\n%s", strings.Join(cmdArgs, " "), string(output))
		return mcp.NewToolResultText(fullResponse), nil
	}
}

type KubectlAnnotateArgs struct {
	Resource      string `json:"resource"`
	Name          string `json:"name,omitempty"`
	Annotations   string `json:"annotations"`
	Namespace     string `json:"namespace,omitempty"`
	LabelSelector string `json:"label_selector,omitempty"`
	Overwrite     bool   `json:"overwrite"`
	All           bool   `json:"all"`
}

func (h *Handler) kubectlAnnotateHandler() mcp.TypedToolHandlerFunc[KubectlAnnotateArgs] {
	return func(
		ctx context.Context,
		req mcp.CallToolRequest,
		args KubectlAnnotateArgs,
	) (*mcp.CallToolResult, error) {
		cmdArgs := []string{"annotate", args.Resource}

		if args.Name != "" {
			cmdArgs = append(cmdArgs, args.Name)
		}

		// Add annotations to the command
		cmdArgs = append(cmdArgs, args.Annotations)

		if args.Namespace != "" {
			cmdArgs = append(cmdArgs, "-n", args.Namespace)
		}
		if args.LabelSelector != "" {
			cmdArgs = append(cmdArgs, "-l", args.LabelSelector)
		}
		if args.Overwrite {
			cmdArgs = append(cmdArgs, "--overwrite")
		}
		if args.All {
			cmdArgs = append(cmdArgs, "--all")
		}

		output, err := h.runKubectl(ctx, cmdArgs...)
		if err != nil {
			errorMsg := fmt.Sprintf("kubectl annotate failed: %v\nCommand: kubectl %s\nOutput: %s",
				err, strings.Join(cmdArgs, " "), string(output))
			return mcp.NewToolResultError(errorMsg), nil
		}

		fullResponse := fmt.Sprintf("Command executed: kubectl %s\n\n%s", strings.Join(cmdArgs, " "), string(output))
		return mcp.NewToolResultText(fullResponse), nil
	}
}
