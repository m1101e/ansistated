package cli

import (
	"fmt"

	"github.com/m1101e/ansistated/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCmd(cfg *config.Config) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "ansistated",
		Short:   "Templatizes pongo2 templates and provides a mechanism for changing variable values",
		Version: "0.1.0",
		PersistentPreRunE: func(rootCmd *cobra.Command, args []string) error {
			viper.BindPFlag("vars_dir", rootCmd.Flags().Lookup("vars-dir"))
			viper.SetDefault("vars_dir", "~/.config/ansistated/vars")

			viper.BindPFlag("templates_dir", rootCmd.Flags().Lookup("templates-dir"))
			viper.SetDefault("templates_dir", "~/.config/ansistated/templates")

			if err := config.LoadConfig(cfg); err != nil {
				return fmt.Errorf("error reading the ansistated configuration file: \n%w", err)
			}

			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd.PersistentFlags().String("vars-dir", "", "A directory containing YAML files with variables")
	rootCmd.PersistentFlags().String("templates-dir", "", "A directory containing template files")

	renderCmd := NewRenderCmd(cfg)
	rootCmd.AddCommand(renderCmd)

	updateCmd := NewUpdateCmd(cfg)
	rootCmd.AddCommand(updateCmd)

	return rootCmd
}
