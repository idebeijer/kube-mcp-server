package cmd

import (
	"fmt"
	"os"

	"github.com/idebeijer/kube-mcp-server/internal/config"
	"github.com/idebeijer/kube-mcp-server/internal/mcpserver"
	"github.com/idebeijer/kube-mcp-server/pkg/logger"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	cfg     *config.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kube-mcp-server",
	Short: "A Kubernetes MCP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := mcpserver.New(cfg, mcpserver.WithResources(), mcpserver.WithTools())
		if err != nil {
			return fmt.Errorf("failed to create MCP server: %w", err)
		}

		switch config.Mode(cfg.Mode) {
		case config.ModeStdio:
			if err := server.StartStdio(); err != nil {
				return fmt.Errorf("failed to start MCP server: %w", err)
			}
		case config.ModeSSE:
			if err := server.StartSSE(fmt.Sprintf(":%s", cfg.SSEPort)); err != nil {
				return fmt.Errorf("failed to start MCP server: %w", err)
			}
		default:
			return fmt.Errorf("unknown mode: %s", cfg.Mode)
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initLogging)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kube-mcp-server.yaml)")

	rootCmd.PersistentFlags().String("log-level", "info", "log level")
	_ = viper.BindPFlag("logLevel", rootCmd.PersistentFlags().Lookup("log-level"))

	rootCmd.PersistentFlags().String("kubeconfig", "", "path to kubeconfig file")
	_ = viper.BindPFlag("kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig"))

	rootCmd.PersistentFlags().String("mode", "stdio", "mode of operation (stdio or sse)")
	_ = viper.BindPFlag("mode", rootCmd.PersistentFlags().Lookup("mode"))

	rootCmd.PersistentFlags().String("sse-port", "8080", "port for SSE mode")
	_ = viper.BindPFlag("ssePort", rootCmd.PersistentFlags().Lookup("sse-port"))
}

func initConfig() {
	var err error
	cfg, err = config.Load(cfgFile)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	cfg.LogLevel = viper.GetString("logLevel")
}

func initLogging() {
	logger.Init(os.Stdout, cfg.LogLevel, false)
}
