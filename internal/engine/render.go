package engine

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/m1101e/ansistated/internal/config"
	"github.com/m1101e/ansistated/internal/minifunc"

	"github.com/flosch/pongo2"
)

func RenderFiles(cfg *config.Config, groups map[string]struct{}) error {
	pongoVars, err := LoadPongoContext(cfg.VarsDir)
	if err != nil {
		return fmt.Errorf("error retrieving pongo2 variables: \n%w", err)
	}

	templateDir, err := minifunc.ExpandTilde(cfg.TemplatesDir)
	if err != nil {
		return err
	}

	_, err = os.Stat(templateDir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("the \"%s\" templates folder does not exist: \n%w", cfg.TemplatesDir, err)
		}
		return err
	}

	mainStream(templateDir, pongoVars, groups)
	return nil
}

func mainStream(templatesDir string, pongoVars *pongo2.Context, groups map[string]struct{}) {
	tasks := make(chan string, 100)
	var wg sync.WaitGroup

	postCmds := renderTemplates(tasks, &wg, pongoVars, groups)

	err := filepath.WalkDir(templatesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("WARNING: Error while traversing the \"%s\" template directory for the \"%s\" file: \n%v", templatesDir, path, err)
			return nil
		}

		if d.IsDir() {
			return nil
		}

		tasks <- path
		return nil
	})

	close(tasks)
	wg.Wait()

	if err != nil {
		log.Printf("WARNING: Error while iterating through the \"%s\" template directory: \n%v", templatesDir, err)
	}

	runPostCmds(*postCmds)
}

func runPostCmds(postCmds map[string]bool) {
	if len(postCmds) > 0 {
		for cmdStr := range postCmds {
			cmd := exec.Command("sh", "-c", cmdStr)
			out, err := cmd.CombinedOutput()
			if err != nil {
				log.Printf("Post command execution error '%s': %v\nOutput: %s\n", cmdStr, err, string(out))
			}
		}
	}
}

func renderTemplates(tasks chan string, wg *sync.WaitGroup, pongoVars *pongo2.Context, groups map[string]struct{}) *map[string]bool {
	var cmdMu sync.Mutex

	numWorkers := runtime.NumCPU()

	postCmds := make(map[string]bool)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for filePath := range tasks {
				data, err := os.ReadFile(filePath)
				if err != nil {
					log.Printf("WARNING: Unable to read the file \"%s\" for templating: \n%v", filePath, err)
					continue
				}

				meta, _, bodyTemplate, err := ParseFile(data)
				if err != nil {
					log.Printf("WARNING: Unable to parse the file data for the \"%s\" template: \n%v", filePath, err)
					continue
				}

				if meta.Dest == "" {
					log.Printf("WARNING: Template \"%s\" has no 'dest' in frontmatter. Skipping.", filePath)
					continue
				}

				matchGroup := false
				if _, exist := groups["all"]; exist {
					matchGroup = true
				} else {
					for reqGroup := range groups {
						if _, exist := meta.GroupsMap[reqGroup]; exist {
							matchGroup = true
							break
						}
					}
				}

				if !matchGroup {
					continue
				}

				pongoTemplate, err := pongo2.FromString(bodyTemplate)
				if err != nil {
					log.Printf("WARNING: Error converting the body of the \"%s\" file template into data for pongo2: \n%v", filePath, err)
					continue
				}

				templatedData, err := pongoTemplate.Execute(*pongoVars)
				if err != nil {
					log.Printf("WARNING: Error while templating the file \"%s\": \n%v", filePath, err)
					continue
				}

				destPath, err := minifunc.ExpandTilde(meta.Dest)
				if err != nil {
					log.Printf("WARNING: Error converting the path \"%s\" which contains a tilde, while applying a template to the file \"%s\": \n%v", meta.Dest, filePath, err)
					continue
				}

				renderedBytes := []byte(templatedData)

				existingBytes, err := os.ReadFile(destPath)
				if err == nil && bytes.Equal(existingBytes, renderedBytes) {
					continue
				}

				if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
					log.Printf("WARNING: Error creating directories for \"%s\": \n%v", destPath, err)
					continue
				}

				err = os.WriteFile(destPath, renderedBytes, 0644)
				if err != nil {
					log.Printf("WARNING: Error writing template to \"%s\": \n%v", destPath, err)
					continue
				}

				if meta.PostCmd != "" {
					cmdMu.Lock()
					postCmds[meta.PostCmd] = true
					cmdMu.Unlock()
				}
			}
		}()
	}

	return &postCmds
}
