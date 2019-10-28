// Package launchpad file functionality.go contains all functionality for output generation.
package launchpad

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// directoryOwner interface allows implementers to specify directory creation.
type directoryOwner interface {
	directoryProperty() *directoryProperty
}

// directoryProperty defines directory to be created.
type directoryProperty struct {
	basename string // Directory name
	dirname  string // ParentId directory name
	backup   bool   // Backup directory during re-creation
}

// path generates the full path of the directory.
func (d *directoryProperty) path() string { return filepath.Join(d.dirname, d.basename) }

// directoryPropertyBackup sets directoryProperty backup property.
func directoryPropertyBackup(backup bool) func(*directoryProperty) error {
	return func(c *directoryProperty) error { c.backup = backup; return nil }
}

// directoryPropertyDirname sets directoryProperty dirname property.
func directoryPropertyDirname(dirname string) func(*directoryProperty) error {
	return func(c *directoryProperty) error { c.dirname = dirname; return nil }
}

// newDirectoryProperty initializes directoryProperty defaulting to current directory with backup turned on.
//
// newDirectoryProperty allows users to specify setter functions to modify output properties.
func newDirectoryProperty(dirname string, options ...func(*directoryProperty) error) *directoryProperty {
	c := &directoryProperty{
		basename: dirname,
		dirname:  ".", // default to binary execution directory
		backup:   true,
	}
	for _, op := range options {
		if err := op(c); err != nil {
			log.Fatalln("Unable to apply directory property modifier", err)
		}
	}
	return c
}

// withDirectory actions on directoryOwner component and creates directory based on directoryProperty.
//
// Backup turned on will rename existing directory with last modify time as postfix.
func withDirectory(comp component) {
	do, ok := comp.(directoryOwner)
	if !ok {
		return
	}
	fp, backup := do.directoryProperty().path(), do.directoryProperty().backup

	if fileInfo, err := os.Stat(fp); err == nil {
		if !backup {
			log.Println("Removing existing folder", fp)
			err = os.RemoveAll(fp)
		} else {
			err = os.Rename(fp, fmt.Sprintf("%s_%d", fp, fileInfo.ModTime().Unix()))
		}

		if err != nil {
			log.Fatalf("Failed to remove/backup existing folder %s", fp)
		}
	}
	log.Printf("Creating directory `%s`", fp)
	if err := os.MkdirAll(fp, 0755); err != nil {
		log.Fatalf("Failed to create folder %s\n", fp)
	}
}

// filesOwner interface allows implementers to specify files creation.
type filesOwner interface {
	files() []file
}

// file interface allows implementers to specify operations for a given file type.
type file interface {
	path() string
	render() string
}

// withFiles processes filesOwner components to creates files.
func withFiles(comp component) {
	fo, ok := comp.(filesOwner)
	if !ok {
		return
	}
	for _, f := range fo.files() {
		if err := ioutil.WriteFile(f.path(), []byte(f.render()), 0644); err != nil {
			log.Fatalln("Failed to generate output", f.path(), err)
		}
	}
}
