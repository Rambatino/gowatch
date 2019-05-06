package app

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

type Watcher interface {
	WatchAndRun() chan error
}

type Watch struct {
	modTimes   []time.Time
	hasChanged chan bool
	basePath   string
	extensions []string
	paths      []string
	recursive  bool
	ignore     []string
	args       []string
	pid        int
}

func NewWatcher(extensions, paths []string, recursive bool, ignore, args []string) (Watcher, error) {
	basePath, err := os.Getwd()
	if err != nil {
		return &Watch{}, err
	}

	return &Watch{[]time.Time{}, make(chan bool), basePath, extensions, paths, recursive, ignore, args, 0}, nil
}

// Watch watches for changes given set of parameters. If extensions passed, will
// look at only those file types (recursively too if passed)
// paths pass will only look in those folders and files (recursively too if passed)
// ignore will ignore all matching folders and files (recursively too if passed)
func (w *Watch) WatchAndRun() chan error {
	go func() {
		for {
			go w.watch()
			time.Sleep(300 * time.Millisecond)
		}
	}()
	for {
		select {
		case changed := <-w.hasChanged:
			if changed {
				w.run()
			}
		}
	}
}

func (w *Watch) run() {
	if w.pid != 0 {
		proc, err := os.FindProcess(w.pid)
		if err != nil {
			log.Println(err)
		}
		// Kill the process
		proc.Kill()
	}
	cmd := exec.Command(w.args[0], w.args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
	w.pid = cmd.Process.Pid
}

func (w *Watch) watch() {
	files, _ := files(w.basePath, w.extensions, w.paths, w.recursive, w.ignore)
	fileTimes := []time.Time{}
	for _, f := range files {
		fileTimes = append(fileTimes, f.ModTime())
	}
	if !reflect.DeepEqual(w.modTimes, fileTimes) {
		w.modTimes = fileTimes
		w.hasChanged <- true
	}
}

func files(basePath string, extensions, paths []string, recursive bool, ignore []string) ([]os.FileInfo, error) {
	files := []os.FileInfo{}
	cleanedPath := filepath.Clean(basePath) + "/"

	filepath.Walk(basePath, func(walkedPath string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			relPath := strings.Replace(walkedPath, cleanedPath, "", 1)
			isSubDir := strings.ContainsAny(relPath, "/")
			appended := false
			for _, p := range paths {
				if strings.HasPrefix(relPath, p) {
					// if there are extensions use them, if not, add it
					if len(extensions) == 0 {
						appended = true
						files = append(files, info)
						continue
					}

					for _, e := range extensions {
						if strings.Contains(info.Name(), "."+e) {
							appended = true
							files = append(files, info)
							continue
						}
					}
				}
			}
			if (!isSubDir || recursive) && !appended {
				if len(paths)+len(extensions) == 0 {
					files = append(files, info)
				} else if matchedPath(relPath, info, extensions) {
					files = append(files, info)
				}
			}
		}
		return nil
	})
	return files, nil
}

func matchedPath(relPath string, fileInfo os.FileInfo, extensions []string) bool {
	for _, e := range extensions {
		if strings.Contains(fileInfo.Name(), "."+e) {
			return true
		}
	}
	return false
}
