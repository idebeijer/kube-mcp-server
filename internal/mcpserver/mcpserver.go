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

	enableTools     bool
	enableResources bool
}

type Option func(*Server)

func WithTools() Option {
	return func(s *Server) {
		s.enableTools = true
	}
}

func WithResources() Option {
	return func(s *Server) {
		s.enableResources = true
	}
}

func New(cfg *config.Config, opts ...Option) (*Server, error) {
	var restCfg *rest.Config
	var err error
	if cfg.Kubeconfig != "" {
		restCfg, err = clientcmd.BuildConfigFromFlags("", cfg.Kubeconfig)
		if err != nil {
			return nil, err
		}
	} else {
		restCfg, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}
	client, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return nil, err
	}

	s := &Server{
		client: client,
	}

	for _, opt := range opts {
		opt(s)
	}

	mcpServerOpts := []server.ServerOption{
		server.WithLogging(),
	}
	if s.enableTools {
		log.Info().Msg("Enabling tools")
		mcpServerOpts = append(mcpServerOpts, server.WithToolCapabilities(true))
	}
	if s.enableResources {
		log.Info().Msg("Enabling resources")
		mcpServerOpts = append(mcpServerOpts, server.WithResourceCapabilities(false, true))
	}

	mcpServer := server.NewMCPServer(
		"kube-mcp-server", "0.1.0",
		mcpServerOpts...,
	)
	s.mcp = mcpServer

	if s.enableTools {
		var toolOpts []tool.Option
		if !cfg.DisableKubectl {
			toolOpts = append(toolOpts, tool.WithKubectlTools())
		}
		tools, _ := tool.NewHandler(s.client, cfg.Kubeconfig, toolOpts...)
		tools.Register(s.mcp)
	}
	if s.enableResources {
		resources := resource.NewHandler(s.client)
		resources.Register(s.mcp)
	}

	return s, nil
}

func (s *Server) StartSSE(addr string) error {
	sse := server.NewSSEServer(s.mcp)
	log.Info().Msgf("Starting MCP server on %s", addr)
	return sse.Start(addr)
}

func (s *Server) StartStdio() error {
	log.Info().Msg("Running in stdio mode. Press Ctrl+C to exit.")
	return server.ServeStdio(s.mcp)
}
