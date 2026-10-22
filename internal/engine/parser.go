package engine

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type FileMeta struct {
	Group     string              `yaml:"group"`
	Dest      string              `yaml:"dest"`
	PostCmd   string              `yaml:"post_cmd"`
	GroupsMap map[string]struct{} `yaml:"-"`
}

func ParseFile(fileData []byte) (*FileMeta, string, string, error) {
	normalizedData := bytes.ReplaceAll(fileData, []byte("\r\n"), []byte("\n"))

	meta := &FileMeta{
		GroupsMap: make(map[string]struct{}),
	}

	if bytes.HasPrefix(normalizedData, []byte("---\n")) {
		parts := bytes.SplitN(normalizedData, []byte("---\n"), 3)

		if len(parts) >= 3 {
			yamlHeader := string(parts[1])
			fileBody := string(parts[2])

			if err := yaml.Unmarshal([]byte(yamlHeader), meta); err != nil {
				return nil, "", "", fmt.Errorf("YAML parsing error in the header: \n%w", err)
			}

			parseGroups(meta)

			return meta, yamlHeader, fileBody, nil
		}
	}

	parseGroups(meta)
	return meta, "", string(normalizedData), nil
}

func parseGroups(meta *FileMeta) {
	meta.GroupsMap = make(map[string]struct{})

	if meta.Group == "" {
		meta.GroupsMap["default"] = struct{}{}
		return
	}

	groups := strings.Split(meta.Group, ",")
	for _, g := range groups {
		cleanGroup := strings.TrimSpace(g)
		if cleanGroup != "" {
			meta.GroupsMap[cleanGroup] = struct{}{}
		}
	}
}
