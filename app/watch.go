package app

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"syscall"
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
	cmd        *exec.Cmd
}

func NewWatcher(extensions, paths []string, recursive bool, ignore, args []string) (Watcher, error) {
	basePath, err := os.Getwd()
	if err != nil {
		return &Watch{}, err
	}

	return &Watch{[]time.Time{}, make(chan bool), basePath, extensions, paths, recursive, ignore, args, nil}, nil
}

// Watch watches for changes given set of parameters. If extensions passed, will
// look at only those file types (recursively too if passed)
// paths pass will only look in those folders and files (recursively too if passed)
// ignore will ignore all matching folders and files (recursively too if passed)
func (w *Watch) WatchAndRun() chan error {
	go func() {
		for {
			go w.watch()
			time.Sleep(1000 * time.Millisecond)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		w.terminate()
		os.Exit(1)
	}()
	for {
		select {
		case changed := <-w.hasChanged:
			if changed {
				w.run()
				fmt.Println("Running:", strings.Join(w.args, " "), " { pid:", w.cmd.Process.Pid, ", fileCount: "+strconv.Itoa(len(w.modTimes))+" }")
			}
		}
	}
}

func (w *Watch) terminate() {
	if w.cmd != nil {
		pgid, err := syscall.Getpgid(w.cmd.Process.Pid)
		if err == nil {
			syscall.Kill(-pgid, 15) // note the minus sign
		}

		w.cmd.Wait()
	}
}

func (w *Watch) run() {
	w.terminate()
	cmd := exec.Command(w.args[0], w.args[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	w.cmd = cmd
	cmd.Start()
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
