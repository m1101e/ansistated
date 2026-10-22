package minifunc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ExpandTilde(path string) (string, error) {
	if path == "~" {
		return os.UserHomeDir()
	}
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("unable to obtain the home directory: \n%w", err)
		}
		return filepath.Join(homeDir, path[2:]), nil
	}
	return path, nil
}
