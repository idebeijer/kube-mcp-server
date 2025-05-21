package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	LogLevel          string `mapstructure:"logLevel"`
	StructuredLogging bool   `mapstructure:"structuredLogging"`
	KubeConfigPath    string `mapstructure:"kubeconfigPath"`
}

func Load(cfgFile string) (*Config, error) {
	v := viper.GetViper()

	// Set default values
	v.SetDefault("logLevel", "info")

	// Configure file-based settings
	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to find home directory: %w", err)
		}
		v.AddConfigPath(home)
		v.SetConfigType("yaml")
		v.SetConfigName(".kube-mcp-server")
	}

	// Read environment variables
	v.AutomaticEnv()

	// Read the config file if it exists
	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Unmarshal into the Config struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
