package engine

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func UpdateVariables(store *YamlStore, target string, isGroup bool, updates map[string]interface{}) error {
	var targetFiles []*YamlFile

	if isGroup {
		files, exists := store.ByGroup[target]
		if !exists {
			return fmt.Errorf("group \"%s\" not found in variables store", target)
		}
		targetFiles = files
	} else {
		file, exists := store.ByPath[target]
		if !exists {
			return fmt.Errorf("file path \"%s\" not found in variables store", target)
		}
		targetFiles = []*YamlFile{file}
	}

	globalMatchedKeys := make(map[string]bool)

	for _, yf := range targetFiles {
		changed := applyUpdatesToNode(yf.Node, updates, yf.Path, "", globalMatchedKeys)

		if changed {
			var buf bytes.Buffer

			if yf.RawHeader != "" {
				buf.WriteString("---\n")
				buf.WriteString(yf.RawHeader)
				buf.WriteString("---\n")
			}

			encoder := yaml.NewEncoder(&buf)
			encoder.SetIndent(2)
			if err := encoder.Encode(yf.Node); err != nil {
				log.Printf("WARNING: Failed to encode updated YAML for \"%s\": \n%v", yf.Path, err)
				continue
			}
			encoder.Close()

			if err := os.WriteFile(yf.Path, buf.Bytes(), 0644); err != nil {
				log.Printf("WARNING: Failed to write updated file \"%s\": \n%v", yf.Path, err)
			} else {
				fmt.Printf("Successfully updated variables in \"%s\"\n", yf.Path)
			}
		}
	}

	warnUnmatchedKeys(updates, globalMatchedKeys, "")

	return nil
}

func warnUnmatchedKeys(updates map[string]interface{}, matched map[string]bool, prefix string) {
	for k, v := range updates {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + "." + k
		}

		if !matched[fullKey] {
			log.Printf("WARNING: Key \"%s\" not found in any targeted file. Skipping update for this key.", fullKey)
		}

		if subUpdates, isMap := v.(map[string]interface{}); isMap {
			warnUnmatchedKeys(subUpdates, matched, fullKey)
		}
	}
}

func applyUpdatesToNode(node *yaml.Node, updates map[string]interface{}, filePath string, prefix string, matchedKeys map[string]bool) bool {
	changed := false

	if node.Kind == yaml.DocumentNode {
		if len(node.Content) > 0 {
			return applyUpdatesToNode(node.Content[0], updates, filePath, prefix, matchedKeys)
		}
		return false
	}

	if node.Kind != yaml.MappingNode {
		log.Printf("WARNING: Expected a mapping node in \"%s\", but got kind %v", filePath, node.Kind)
		return false
	}

	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valNode := node.Content[i+1]

		fullKey := keyNode.Value
		if prefix != "" {
			fullKey = prefix + "." + keyNode.Value
		}

		updateVal, exists := updates[keyNode.Value]
		if !exists {
			continue
		}

		matchedKeys[fullKey] = true

		if subUpdates, isMap := updateVal.(map[string]interface{}); isMap {
			if valNode.Kind == yaml.MappingNode {
				if applyUpdatesToNode(valNode, subUpdates, filePath, fullKey, matchedKeys) {
					changed = true
				}
			} else {
				log.Printf("WARNING: Type mismatch for key \"%s\" in \"%s\". Expected nested mapping.", fullKey, filePath)
			}
		} else {
			var tmpNode yaml.Node
			if err := tmpNode.Encode(updateVal); err == nil {
				if valNode.Kind == yaml.ScalarNode && tmpNode.Kind == yaml.ScalarNode {
					if valNode.Value != tmpNode.Value {
						valNode.Value = tmpNode.Value
						valNode.Tag = tmpNode.Tag
						changed = true
					}
				} else {
					tmpNode.HeadComment = valNode.HeadComment
					tmpNode.LineComment = valNode.LineComment
					tmpNode.FootComment = valNode.FootComment
					*valNode = tmpNode
					changed = true
				}
			} else {
				log.Printf("WARNING: Failed to encode value for key \"%s\" in \"%s\": %v", fullKey, filePath, err)
			}
		}
	}

	return changed
}
