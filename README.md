# Gowatch

Tell gowatch what type of files you want watching, and then what you want to execute when those files change. There are other options out there, but the aim of this is to be as simple and all encompassing as possible.

## How to install

Via homebrew:

``` bash
brew tap rambatino/gowatch
brew install gowatch
```

or as a go package:

``` bash
go get github.com/Rambatino/gowatch
```

## Usage

To understand what the commands do, it's helpful to use the `-h` flag.

E.g. `gowatch -h`

``` bash
Available Commands:
  help        Help about any command
  run         Run custom command
  version     Print the version number of GoWatch
```

### Examples
``` bash
# rerun your go server on file changes
gowatch run -e=go -r go run *.go

# re-test your typescript react app on changes (don't include node modules)
gowatch run -e=js,ts,tx,jx -o=node_modules -r yarn test-watch

# rerun your ruby script on file change
gowatch run ruby starter.rb # just read all files in top level directory
gowatch run -e=rb ruby starter.rb # only your ruby file
```

### All Flags

To access all the flags and make the most out of the command, run `gowatch run -h`

```
Usage:
  gowatch run [flags]

Flags:
  -e, --extensions strings          Comma separated file extensions/types in which to search for. N.B. don't pass globs.
  -h, --help                        help for run
  -i, --ignore-extensions strings   What file extensions to ignore. N.B. don't pass globs.
  -o, --ignore-paths strings        What file paths to ignore. N.B. don't pass globs.
  -p, --paths strings               Comma separated paths (folders and files) in which to search in. N.B. don't pass globs.
  -r, --recursive                   Whether to search recursively
```

Passing Globs e.g. *.go as arguments to the flags will only result in issues, primarily that shells such as zsh will immediately turn *.go into, say, `main.go server.go` which will confuse the cli as it will only search for those files then.
