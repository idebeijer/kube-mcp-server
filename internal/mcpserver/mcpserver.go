package mcpserver

import (
	"github.com/idebeijer/kube-mcp-server/internal/config"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Server struct {
	MCP *server.MCPServer
	K8s *kubernetes.Clientset
}

func New(cfg *config.Config) (*Server, error) {
	var restCfg *rest.Config
	if cfg.KubeConfigPath != "" {
		restCfg, _ = clientcmd.BuildConfigFromFlags("", cfg.KubeConfigPath)
	} else {
		restCfg, _ = rest.InClusterConfig()
	}
	clientset, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return nil, err
	}

	m := server.NewMCPServer(
		"kube-mcp-server", "0.1.0",
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)

	RegisterTools(m, clientset)

	return &Server{
		MCP: m,
		K8s: clientset,
	}, nil
}

func (s *Server) Start(addr string) error {
	sse := server.NewSSEServer(s.MCP)
	log.Info().Msgf("starting MCP server on %s", addr)
	return sse.Start(addr)
}

func (s *Server) StartStdio() error {
	return server.ServeStdio(s.MCP)
}
