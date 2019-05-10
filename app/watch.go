package app

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

type Watcher interface {
	WatchAndRun() chan error
}

type Watch struct {
	files   *files
	args    []string
	cmd     *exec.Cmd
	watcher *fsnotify.Watcher
}

func NewWatcher(extensions, paths []string, recursive bool, ignoreExtensions, ignorePaths, args []string) (Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return &Watch{}, err
	}

	f, err := NewFiles(extensions, paths, recursive, ignoreExtensions, ignorePaths)
	if err != nil {
		return &Watch{}, err
	}

	w := Watch{args: args, watcher: watcher, files: f}
	go w.listenForExit()

	return &w, nil
}

func (w *Watch) listenForExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	w.terminate()
	os.Exit(1)
}

// Watch watches for changes given set of parameters. If extensions passed, will
// look at only those file types (recursively too if passed)
// paths pass will only look in those folders and files (recursively too if passed)
// ignore will ignore all matching folders and files (recursively too if passed)
func (w *Watch) WatchAndRun() chan error {
	go func() {
		for {
			select {
			case event, ok := <-w.watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					w.run()
				}
			case err, ok := <-w.watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error:", err)
			case <-w.files.changed:
				for _, f := range w.files.foundFiles {
					w.watcher.Add(f)
				}
				w.run()
			}
		}
	}()

	return nil
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

func (w *Watch) terminate() {
	if w.cmd != nil {
		pgid, err := syscall.Getpgid(w.cmd.Process.Pid)
		if err == nil {
			syscall.Kill(-pgid, 15)
		}
		w.cmd.Wait()
	}
}
