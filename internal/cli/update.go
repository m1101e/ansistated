package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/m1101e/ansistated/internal/config"
	"github.com/m1101e/ansistated/internal/engine"

	"github.com/spf13/cobra"
)

func NewUpdateCmd(cfg *config.Config) *cobra.Command {
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Reads JSON from stdin and updates the values of variables in the selected file or group",
		RunE: func(cmd *cobra.Command, args []string) error {
			group, err := cmd.Flags().GetString("group")
			if err != nil {
				return fmt.Errorf("error retrieving the \"group\" flag: \n%w", err)
			}
			file, err := cmd.Flags().GetString("file")
			if err != nil {
				return fmt.Errorf("error retrieving the \"file\" flag: \n%w", err)
			}

			target := ""
			isGroup := false
			if group != "" {
				target = group
				isGroup = true
			} else if file != "" {
				target = file
			} else {
				return fmt.Errorf("you cannot run the \"update\" command without the \"--group\" or \"--file\" flag, or by specifying an empty string as the target")
			}

			inputData, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read from stdin: \n%w", err)
			}

			var updates map[string]interface{}
			if err := json.Unmarshal(inputData, &updates); err != nil {
				return fmt.Errorf("failed to parse JSON input: \n%w", err)
			}

			yamlStore, err := engine.LoadYamlStore(cfg.VarsDir)
			if err != nil {
				return fmt.Errorf("error loading variables: \n%w", err)
			}

			err = engine.UpdateVariables(yamlStore, target, isGroup, updates)
			if err != nil {
				return fmt.Errorf("failed to update variables: \n%w", err)
			}

			fmt.Println("Variables have been updated")

			return nil
		},
	}

	updateCmd.PersistentFlags().StringP("group", "g", "", "Sets the group for updating variables")
	updateCmd.PersistentFlags().StringP("file", "f", "", "Specifies the file (use a relative path from “vars”) to update the variables")

	updateCmd.MarkFlagsMutuallyExclusive("group", "file")

	return updateCmd
}
