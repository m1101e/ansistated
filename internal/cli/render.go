package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/m1101e/ansistated/internal/config"
	"github.com/m1101e/ansistated/internal/engine"
)

func NewRenderCmd(cfg *config.Config) *cobra.Command {
	renderCmd := &cobra.Command{
		Use:   "render",
		Short: "Applies a template to all files in the folder",
		RunE: func(cmd *cobra.Command, args []string) error {
			groups, err := cmd.Flags().GetStringSlice("groups")
			if err != nil {
				return fmt.Errorf("Error retrieving the \"groups\" flag: \n%w", err)
			}
			groupsMap := make(map[string]struct{})
			for _, group := range groups {
				groupsMap[group] = struct{}{}
			}
			err = engine.RenderFiles(cfg, groupsMap)
			if err != nil {
				return err
			}
			return nil
		},
	}
	renderCmd.PersistentFlags().StringSlice("groups", []string{"all"}, "Sets the group for templating")

	return renderCmd
}
