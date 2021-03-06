package app

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

var flagtests = map[string]struct {
	extensions       []string
	paths            []string
	recursive        bool
	ignoreExtensions []string
	ignorePaths      []string
	out              []string
	err              error
}{
	"go paths":                               {[]string{"go"}, []string{}, false, []string{}, []string{}, []string{"lol.go", "main.go"}, nil},
	"go paths recursive":                     {[]string{"go"}, []string{}, true, []string{}, []string{}, []string{"code/app.go", "lol.go", "main.go"}, nil},
	"all paths recursive":                    {[]string{}, []string{}, true, []string{}, []string{}, []string{"cat.sh", "code/app.go", "code/app.js", "lol.go", "main.go", "node_modules/run.js"}, nil},
	"all paths not recursive":                {[]string{}, []string{}, false, []string{}, []string{}, []string{"cat.sh", "lol.go", "main.go"}, nil},
	"only main.go":                           {[]string{}, []string{"main.go"}, false, []string{}, []string{}, []string{"main.go"}, nil},
	"main.go and cat.sh":                     {[]string{}, []string{"main.go", "cat.sh"}, false, []string{}, []string{}, []string{"cat.sh", "main.go"}, nil},
	"main.go and code/app.go":                {[]string{}, []string{"main.go", "code/app.go"}, false, []string{}, []string{}, []string{"code/app.go", "main.go"}, nil},
	"main.go and code":                       {[]string{}, []string{"main.go", "code"}, false, []string{}, []string{}, []string{"code/app.go", "code/app.js", "main.go"}, nil},
	"code and only go files in folder: code": {[]string{"go"}, []string{"main.go", "code"}, false, []string{}, []string{}, []string{"code/app.go", "lol.go", "main.go"}, nil},
	"only node_modules":                      {[]string{}, []string{"node_modules"}, false, []string{}, []string{}, []string{"node_modules/run.js"}, nil},
	"node_modules and code":                  {[]string{}, []string{"node_modules", "code"}, false, []string{}, []string{}, []string{"code/app.go", "code/app.js", "node_modules/run.js"}, nil},
	"dodgy":                                  {[]string{}, []string{"nodesad../asd&&&***_modules", "code"}, false, []string{}, []string{}, []string{"code/app.go", "code/app.js"}, nil},
	"all paths recursive ignore js":          {[]string{}, []string{}, true, []string{"js"}, []string{}, []string{"cat.sh", "code/app.go", "lol.go", "main.go"}, nil},
	"all paths recursive ignore code dir":    {[]string{}, []string{}, true, []string{}, []string{"code"}, []string{"cat.sh", "lol.go", "main.go", "node_modules/run.js"}, nil},
}

func TestFiles(t *testing.T) {
	for key, tt := range flagtests {
		t.Run(key, func(t *testing.T) {
			basePath := AddFoldersAndFiles()
			f := files{
				extensions:       tt.extensions,
				paths:            tt.paths,
				recursive:        tt.recursive,
				ignoreExtensions: tt.ignoreExtensions,
				ignorePaths:      tt.ignorePaths,
				basePath:         AddFoldersAndFiles(),
			}
			ff, err := f.findFiles()
			ffSpliced := []string{}
			for _, p := range ff {
				ffSpliced = append(ffSpliced, strings.Replace(p, basePath+"/", "", 1))
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
