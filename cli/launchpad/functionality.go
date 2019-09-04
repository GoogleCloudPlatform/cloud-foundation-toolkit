// Package launchpad file functionality.go contains all functionality for output generation.
package launchpad

import (
	"errors"
	"fmt"
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

// path generates the full path of the directory
func (d *directoryProperty) path() string { return filepath.Join(d.dirname, d.basename) }

// directoryPropertyBackup sets directoryProperty backup property.
func directoryPropertyBackup(backup bool) func(*directoryProperty) error {
	return func(c *directoryProperty) error { c.backup = backup; return nil }
}

// directoryPropertyDirname sets directoryProperty dirname property.
func directoryPropertyDirname(dirname string) func(*directoryProperty) error {
	return func(c *directoryProperty) error { c.dirname = dirname; return nil }
}

// newDirectoryProperty initializes directoryProperty defaulting to current directly with backup.
//
// newDirectoryProperty allows users to specify setter functions to modify default output.
func newDirectoryProperty(dirname string, options ...func(*directoryProperty) error) *directoryProperty {
	c := &directoryProperty{
		basename: dirname,
		dirname:  ".", // default to binary execution directory
		backup:   true,
	}
	for _, op := range options {
		err := op(c)
		if err != nil {
			log.Fatalln("Unable to process directory property", err)
		}
	}
	return c
}

// withDirectory actions on components implements directoryOwner interface and creates directory as specified.
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
	err := os.MkdirAll(fp, 0755)
	if err != nil {
		log.Fatalf("Failed to create folder %s\n", fp)
	}
}

// filesOwner interface allows implementers to specify file creation.
type filesOwner interface {
	files() []file
}

// file interface allows implementers to specify specific file operations for their type.
type file interface {
	path() string
	render() string
}

// withFiles actions on components implements filesOwner interface and creates files as specified.
func withFiles(comp component) {
	fo, ok := comp.(filesOwner)
	if !ok {
		return
	}
	for _, f := range fo.files() {
		err := writeFile(f.path(), f.render())
		if err != nil {
			log.Fatalln("Failed to generate output", f.path(), err)
		}
	}
}

// writeFile creates or replace a file based path and content provided.
func writeFile(fp string, content string) error {
	if _, err := os.Stat(fp); err == nil {
		err := os.Remove(fp)
		if err != nil {
			return errors.New(fmt.Sprintln("Unable to remove file", fp))
		}
	}
	fd, err := os.Create(fp)
	if err != nil {
		return errors.New(fmt.Sprintln("Unable to create file", fp))
	}
	_, err = fd.WriteString(content)
	if err != nil {
		return errors.New(fmt.Sprintln("Unable to write to file", fp))
	}
	err = fd.Close()
	if err != nil {
		return errors.New(fmt.Sprintln("Unable to close file", fp))
	}
	return nil
}
