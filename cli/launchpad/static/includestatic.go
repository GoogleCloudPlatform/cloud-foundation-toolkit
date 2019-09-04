package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Process all files under launchpad/static/**/*.tmpl and encode
// them as string literals
func main() {
	fd, err := os.Create(targetFile("statics.go"))
	if err != nil {
		panic(err)
	}
	write(fd, []byte("package launchpad\n// WARNING: Generated file, do not modify directly!\n\nvar statics = map[string]string {\n"))
	fsWalk(fd, targetFile("static"))
	write(fd, []byte("}\n"))
	err = fd.Close()
	if err != nil {
		panic(err)
	}
}

// Getting relative path of cli/launchpad/statics.go
func targetFile(file string) string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if filepath.Base(dir) == "launchpad" {
		return fmt.Sprintf("./%s", file)
	} else if filepath.Base(dir) == "cli" {
		return fmt.Sprintf("launchpad/%s", file)
	} else {
		panic(errors.New("unrecognized file location"))
	}
}

// Recursively collect *.tmpl files from the given directory and subdirectory and output as string literal
func fsWalk(fd *os.File, dir string) {
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, f := range fs {
		fp := filepath.Join(dir, f.Name())
		if f.IsDir() {
			fsWalk(fd, fp)
		}
		if !(strings.HasSuffix(f.Name(), ".tmpl") || strings.HasSuffix(f.Name(), ".yaml")) {
			continue
		}
		write(fd, []byte(fmt.Sprintf("\t\"%s\": `", fp)))
		fdFile, err := os.Open(fp)
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(fd, fdFile)
		if err != nil {
			panic(err)
		}
		err = fdFile.Close()
		if err != nil {
			panic(err)
		}
		write(fd, []byte("`,\n"))
	}
}

func write(fd *os.File, str []byte) {
	_, err := fd.Write(str)
	if err != nil {
		panic(err)
	}
}
