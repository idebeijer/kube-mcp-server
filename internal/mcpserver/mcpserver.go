package mcpserver

import (
	"github.com/idebeijer/kube-mcp-server/internal/config"
	"github.com/idebeijer/kube-mcp-server/pkg/resource"
	"github.com/idebeijer/kube-mcp-server/pkg/tool"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Server struct {
	mcp    *server.MCPServer
	client *kubernetes.Clientset
}

type Option func(*Server)

func WithTools() Option {
	return func(k *Server) {
		tools := tool.NewHandler(k.client)
		tools.Register(k.mcp)
	}
}

func WithResources() Option {
	return func(s *Server) {
		resources := resource.NewHandler(s.client)
		resources.Register(s.mcp)
	}
}

func New(cfg *config.Config, opts ...Option) (*Server, error) {
	var restCfg *rest.Config
	if cfg.Kubeconfig != "" {
		restCfg, _ = clientcmd.BuildConfigFromFlags("", cfg.Kubeconfig)
	} else {
		restCfg, _ = rest.InClusterConfig()
	}
	client, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return nil, err
	}

	mcpServer := server.NewMCPServer(
		"kube-mcp-server", "0.1.0",
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	s := &Server{
		mcp:    mcpServer,
		client: client,
	}
	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}

func (s *Server) StartSSE(addr string) error {
	sse := server.NewSSEServer(s.mcp)
	log.Info().Msgf("starting MCP server on %s", addr)
	return sse.Start(addr)
}

func (s *Server) StartStdio() error {
	return server.ServeStdio(s.mcp)
}
