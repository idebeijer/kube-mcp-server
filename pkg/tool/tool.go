package tool

import (
	"fmt"
	"os/exec"

	"github.com/mark3labs/mcp-go/server"
	"k8s.io/client-go/kubernetes"
)

type Handler struct {
	client         *kubernetes.Clientset
	kubeconfigPath string

	kubectlEnabled bool
	kubectlPath    string
}

type Option func(handler *Handler)

func WithKubectlTools() Option {
	return func(h *Handler) {
		h.kubectlEnabled = true
	}
}

func NewHandler(client *kubernetes.Clientset, kubeconfigPath string, opts ...Option) (*Handler, error) {
	h := &Handler{
		client:         client,
		kubeconfigPath: kubeconfigPath,
	}
	for _, opt := range opts {
		opt(h)
	}

	if h.kubectlEnabled {
		path, err := exec.LookPath("kubectl")
		if err != nil {
			return nil, fmt.Errorf("kubectl not installed or not found in PATH: %w", err)
		}
		h.kubectlPath = path
	}

	return h, nil
}

func (h *Handler) Register(m *server.MCPServer) {
	h.registerPods(m)

	if h.kubectlEnabled {
		h.registerKubectl(m)
	}
}
