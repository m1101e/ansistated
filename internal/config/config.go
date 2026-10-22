package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type AppFile struct {
	Src     string `mapstructure:"src"`
	Dest    string `mapstructure:"dest"`
	PostCmd string `mapstructure:"post_cmd"`
}

type Config struct {
	VarsDir      string `mapstructure:"vars_dir"`
	TemplatesDir string `mapstructure:"templates_dir"`
}

func LoadConfig(cfg *Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("unable to obtain the home directory: \n%w", err)
	}

	configPath := filepath.Join(homeDir, ".config/ansistated/ansistated.yaml")

	viper.SetConfigFile(configPath)

	_ = viper.ReadInConfig()

	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("failed to Unmarshal data: \n%w", err)
	}

	return nil
}
