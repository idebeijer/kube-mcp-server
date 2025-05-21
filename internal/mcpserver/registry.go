package mcpserver

import (
	"github.com/idebeijer/kube-mcp-server/internal/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"k8s.io/client-go/kubernetes"
)

func RegisterTools(m *server.MCPServer, clientset *kubernetes.Clientset) {
	m.AddTool(tools.CountPodsToolDefinition(), mcp.NewTypedToolHandler[tools.CountPodsArgs](tools.CountPodsHandler(clientset)))
}
