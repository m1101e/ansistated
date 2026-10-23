package cli

import (
	"encoding/json"
	"fmt"

	"github.com/m1101e/ansistated/internal/config"
	"github.com/m1101e/ansistated/internal/engine"
	"github.com/m1101e/ansistated/internal/minifunc"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewReadCmd(cfg *config.Config) *cobra.Command {
	readCmd := &cobra.Command{
		Use:   "read",
		Short: "Outputs merged variable values to stdout in the selected format",
		RunE: func(cmd *cobra.Command, args []string) error {
			group, err := cmd.Flags().GetString("group")
			if err != nil {
				return fmt.Errorf("error retrieving the \"group\" flag: \n%w", err)
			}

			file, err := cmd.Flags().GetString("file")
			if err != nil {
				return fmt.Errorf("error retrieving the \"file\" flag: \n%w", err)
			}

			if cmd.Flags().Changed("file") && file == "" {
				return fmt.Errorf("you cannot set the \"file\" flag and specify an empty string as its value")
			}

			yamlStore, err := engine.LoadYamlStore(cfg.VarsDir)
			if err != nil {
				return fmt.Errorf("error loading variables: \n%w", err)
			}

			groupChanged := cmd.Flags().Changed("group")
			fileChanged := cmd.Flags().Changed("file")

			result := make(map[string]interface{})

			if fileChanged {
				yamlFile, exist := yamlStore.ByPath[file]
				if !exist {
					return fmt.Errorf("file not found: %s", file)
				}

				var current map[string]interface{}
				if err := yamlFile.Node.Decode(&current); err != nil {
					return fmt.Errorf("decoding error: \n%w", err)
				}
				result = current
			} else {
				var files []*engine.YamlFile
				if groupChanged {
					var exist bool
					files, exist = yamlStore.ByGroup[group]
					if !exist {
						return fmt.Errorf("group not found: %s", group)
					}
				} else {
					for _, f := range yamlStore.ByPath {
						files = append(files, f)
					}
				}

				for _, yamlFile := range files {
					var current map[string]interface{}
					if err := yamlFile.Node.Decode(&current); err != nil {
						return fmt.Errorf("decoding error: \n%w", err)
					}
					minifunc.MergeMaps(result, current)
				}
			}

			isYaml, err := cmd.Flags().GetBool("yaml")
			if err != nil {
				return fmt.Errorf("error retrieving the \"yaml\" flag: \n%w", err)
			}

			var encoded []byte
			if isYaml {
				encoded, err = yaml.Marshal(result)
				if err != nil {
					return fmt.Errorf("yaml marshaling error: \n%w", err)
				}
			} else {
				encoded, err = json.MarshalIndent(result, "", "  ")
				if err != nil {
					return fmt.Errorf("json marshaling error: \n%w", err)
				}
			}

			fmt.Println(string(encoded))
			return nil
		},
	}

	readCmd.PersistentFlags().StringP("group", "g", "all", "Sets the group for reading variables")
	readCmd.PersistentFlags().StringP("file", "f", "", "Specifies the file (use a relative path from “vars”) to read the variables")
	readCmd.MarkFlagsMutuallyExclusive("group", "file")

	readCmd.PersistentFlags().Bool("json", true, "Sets the data output format to JSON")
	readCmd.PersistentFlags().Bool("yaml", false, "Sets the data output format to YAML")
	readCmd.MarkFlagsMutuallyExclusive("json", "yaml")

	return readCmd
}
