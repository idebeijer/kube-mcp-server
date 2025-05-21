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
		log.Info().Msg("Starting kube-mcp-server")
		log.Debug().Msgf("Config: %+v", cfg)
		server, err := mcpserver.New(cfg)
		if err != nil {
			return err
		}
		//if err := server.Start(":8088"); err != nil {
		//	return fmt.Errorf("failed to start server: %w", err)
		//}
		if err := server.StartStdio(); err != nil {
			return fmt.Errorf("failed to start server: %w", err)
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

	rootCmd.PersistentFlags().String("kubeconfig-path", "", "path to kubeconfig file")
	_ = viper.BindPFlag("kubeconfigPath", rootCmd.PersistentFlags().Lookup("kubeconfig-path"))
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
