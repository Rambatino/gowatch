package app

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

type files struct {
	foundFiles []string
	changed    chan bool

	basePath         string
	paths            []string
	extensions       []string
	recursive        bool
	ignoreExtensions []string
	ignorePaths      []string
}

func NewFiles(extensions, paths []string, recursive bool, ignoreExtensions, ignorePaths []string) (*files, error) {
	basePath, err := os.Getwd()
	if err != nil {
		return &files{}, err
	}
	f := files{
		basePath:         basePath,
		paths:            paths,
		extensions:       extensions,
		recursive:        recursive,
		ignoreExtensions: ignoreExtensions,
		ignorePaths:      ignorePaths,
		changed:          make(chan bool),
	}
	go func() {
		for {
			files, err := f.findFiles()
			if err != nil {
				fmt.Println(err.Error())
			}
			if !reflect.DeepEqual(files, f.foundFiles) {
				f.foundFiles = files
				f.changed <- true
			}
			time.Sleep(5 * time.Second)
		}
	}()

	return &f, nil
}

func (f *files) findFiles() ([]string, error) {
	files := []string{}
	cleanedPath := filepath.Clean(f.basePath) + "/"

	filepath.Walk(f.basePath, func(walkedPath string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			relPath := strings.Replace(walkedPath, cleanedPath, "", 1)
			if matchedPath(relPath, f.ignoreExtensions, f.ignorePaths) {
				return nil
			}

			isSubDir := strings.ContainsAny(relPath, "/")
			appended := false
			for _, p := range f.paths {
				if strings.HasPrefix(relPath, p) {
					// if there are extensions use them, if not, add it
					if len(f.extensions) == 0 {
						appended = true
						files = append(files, walkedPath)
						continue
					}

					for _, e := range f.extensions {
						if strings.HasSuffix(info.Name(), "."+e) {
							appended = true
							files = append(files, walkedPath)
							continue
						}
					}
				}
			}
			if (!isSubDir || f.recursive) && !appended {
				if len(f.paths)+len(f.extensions) == 0 {
					files = append(files, walkedPath)
				} else if matchedPath(relPath, f.extensions, f.paths) {
					files = append(files, walkedPath)
				}
			}
		}
		return nil
	})
	return files, nil
}

func matchedPath(relPath string, extensions, paths []string) bool {
	for _, e := range extensions {
		if strings.HasSuffix(relPath, "."+e) {
			return true
		}
	}
	for _, p := range paths {
		if strings.HasPrefix(relPath, p) {
			return true
		}
	}
	return false
}
