package engine

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/m1101e/ansistated/internal/minifunc"

	"github.com/flosch/pongo2"
	"gopkg.in/yaml.v3"
)

func LoadPongoContext(dirPath string) (*pongo2.Context, error) {
	expandedPath, err := minifunc.ExpandTilde(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error expanding path for variables directory: \n%w", err)
	}

	if _, err := os.Stat(expandedPath); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("the \"%s\" variables folder does not exist: \n%w", expandedPath, err)
		}
		return nil, err
	}

	mergedVars := make(map[string]interface{})
	var mu sync.Mutex
	var wg sync.WaitGroup
	tasks := make(chan string, 100)

	numWorkers := runtime.NumCPU()
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range tasks {
				fileData, err := os.ReadFile(path)
				if err != nil {
					log.Printf("WARNING: Failed to read file \"%s\": \n%v", path, err)
					continue
				}

				_, _, body, err := ParseFile(fileData)
				if err != nil {
					log.Printf("WARNING: Failed to parse file \"%s\": \n%v", path, err)
					continue
				}

				var fileVars map[string]interface{}
				if err := yaml.Unmarshal([]byte(body), &fileVars); err != nil {
					log.Printf("WARNING: YAML unmarshal error in \"%s\": \n%v", path, err)
					continue
				}

				mu.Lock()
				minifunc.MergeMaps(mergedVars, fileVars)
				mu.Unlock()
			}
		}()
	}

	err = filepath.WalkDir(expandedPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("WARNING: Error accessing path \"%s\": \n%v", path, err)
			return nil
		}
		if !d.IsDir() && (filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml") {
			tasks <- path
		}
		return nil
	})

	close(tasks)
	wg.Wait()

	if err != nil {
		return nil, fmt.Errorf("error while walking variables directory: \n%w", err)
	}

	pongoCtx := pongo2.Context(mergedVars)
	return &pongoCtx, nil
}

type YamlFile struct {
	Path      string
	Meta      *FileMeta
	RawHeader string
	Node      *yaml.Node
}

type YamlStore struct {
	ByPath  map[string]*YamlFile
	ByGroup map[string][]*YamlFile
}

func LoadYamlStore(dirPath string) (*YamlStore, error) {
	expandedPath, err := minifunc.ExpandTilde(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error expanding path for variables directory: \n%w", err)
	}

	store := &YamlStore{
		ByPath:  make(map[string]*YamlFile),
		ByGroup: make(map[string][]*YamlFile),
	}

	if _, err := os.Stat(expandedPath); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("the \"%s\" variables folder does not exist: \n%w", expandedPath, err)
		}
		return nil, err
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	tasks := make(chan string, 100)

	numWorkers := runtime.NumCPU()
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range tasks {
				fileData, err := os.ReadFile(path)
				if err != nil {
					log.Printf("WARNING: Failed to read file \"%s\": \n%v", path, err)
					continue
				}

				relPath, err := filepath.Rel(expandedPath, path)
				if err != nil {
					log.Printf("WARNING: Failed to determine relative path for \"%s\": \n%v", path, err)
					continue
				}

				meta, rawHeader, body, err := ParseFile(fileData)
				if err != nil {
					log.Printf("WARNING: Failed to parse file \"%s\": \n%v", path, err)
					continue
				}

				yf := &YamlFile{
					Path:      path,
					Meta:      meta,
					RawHeader: rawHeader,
					Node:      &yaml.Node{},
				}
				if err := yaml.Unmarshal([]byte(body), yf.Node); err != nil {
					log.Printf("WARNING: YAML AST unmarshal error in \"%s\": \n%v", path, err)
					continue
				}

				groupName := meta.Group
				if groupName == "" {
					groupName = "default"
				}

				mu.Lock()
				store.ByPath[relPath] = yf

				for groupName := range meta.GroupsMap {
					store.ByGroup[groupName] = append(store.ByGroup[groupName], yf)
				}
				mu.Unlock()
			}
		}()
	}

	err = filepath.WalkDir(expandedPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("WARNING: Error accessing path \"%s\": \n%v", path, err)
			return nil
		}
		if !d.IsDir() && (filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml") {
			tasks <- path
		}
		return nil
	})

	close(tasks)
	wg.Wait()

	if err != nil {
		return nil, fmt.Errorf("error while walking variables directory: \n%w", err)
	}

	return store, nil
}
