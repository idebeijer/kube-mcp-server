package resource

import (
	"github.com/mark3labs/mcp-go/server"
	"k8s.io/client-go/kubernetes"
)

type Handler struct {
	client *kubernetes.Clientset
}

func NewHandler(client *kubernetes.Clientset) *Handler {
	return &Handler{
		client: client,
	}
}

func (h *Handler) Register(m *server.MCPServer) {
	h.registerPods(m)
	h.registerDeployments(m)
	h.registerServices(m)
	h.registerStatefulSets(m)
}
