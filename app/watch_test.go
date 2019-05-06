package app

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

var flagtests = map[string]struct {
	extensions []string
	paths      []string
	recursive  bool
	ignore     []string
	out        []string
	err        error
}{
	"go paths":                               {[]string{"go"}, []string{}, false, []string{}, []string{"lol.go", "main.go"}, nil},
	"go paths recursive":                     {[]string{"go"}, []string{}, true, []string{}, []string{"app.go", "lol.go", "main.go"}, nil},
	"all paths recursive":                    {[]string{}, []string{}, true, []string{}, []string{"cat.sh", "app.go", "app.js", "lol.go", "main.go", "run.js"}, nil},
	"all paths not recursive":                {[]string{}, []string{}, false, []string{}, []string{"cat.sh", "lol.go", "main.go"}, nil},
	"only main.go":                           {[]string{}, []string{"main.go"}, false, []string{}, []string{"main.go"}, nil},
	"main.go and cat.sh":                     {[]string{}, []string{"main.go", "cat.sh"}, false, []string{}, []string{"cat.sh", "main.go"}, nil},
	"main.go and code/app.go":                {[]string{}, []string{"main.go", "code/app.go"}, false, []string{}, []string{"app.go", "main.go"}, nil},
	"main.go and code":                       {[]string{}, []string{"main.go", "code"}, false, []string{}, []string{"app.go", "app.js", "main.go"}, nil},
	"code and only go files in folder: code": {[]string{"go"}, []string{"main.go", "code"}, false, []string{}, []string{"app.go", "lol.go", "main.go"}, nil},
	"only node_modules":                      {[]string{}, []string{"node_modules"}, false, []string{}, []string{"run.js"}, nil},
	"node_modules and code":                  {[]string{}, []string{"node_modules", "code"}, false, []string{}, []string{"app.go", "app.js", "run.js"}, nil},
	"dodgy":                                  {[]string{}, []string{"nodesad../asd&&&***_modules", "code"}, false, []string{}, []string{"app.go", "app.js"}, nil},
}

//:TODO need to test ignores too

func TestFiles(t *testing.T) {
	for key, tt := range flagtests {
		t.Run(key, func(t *testing.T) {
			basePath := AddFoldersAndFiles()
			ff, err := files(AddFoldersAndFiles(), tt.extensions, tt.paths, tt.recursive, tt.ignore)
			ffSpliced := []string{}
			for _, p := range ff {
				ffSpliced = append(ffSpliced, strings.Replace(p.Name(), basePath+"/", "", 1))
			}
			if reflect.DeepEqual(ffSpliced, tt.out) == false {
				t.Errorf("got %q, want %q", ffSpliced, tt.out)
			}
			if err != nil && tt.err == nil || tt.err != nil && err == nil {
				t.Error("errors do not match")
			}
		})
	}
}

func AddFoldersAndFiles() string {
	baseFolder := os.TempDir() + "test"
	os.RemoveAll(baseFolder)
	os.MkdirAll(baseFolder, os.ModePerm)

	foldersToCreate := []string{"code", "node_modules"}
	for _, f := range foldersToCreate {
		os.MkdirAll(baseFolder+"/"+f, os.ModePerm)
	}

	filesToCreate := []string{"main.go", "lol.go", "cat.sh", "node_modules/run.js", "code/app.js", "code/app.go"}
	for _, f := range filesToCreate {
		os.Create(baseFolder + "/" + f)
	}
	return baseFolder
}
