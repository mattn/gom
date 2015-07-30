package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Use a wrapper to differentiate logged panics from unexpected ones.
type LoggedError struct{ error }

func panicOnError(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Abort: %s: %s\n", msg, err)
		panic(err)
	}
}

func errorf(format string, args ...interface{}) {
	// Ensure the user's command prompt starts on the next line.
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(os.Stderr, format, args...)
	panic(LoggedError{}) // Panic instead of os.Exit so that deferred will run.
}

func mustCopyFile(destFilename, srcFilename string) {
	destFile, err := os.Create(destFilename)
	panicOnError(err, "Failed to create file "+destFilename)

	srcFile, err := os.Open(srcFilename)
	panicOnError(err, "Failed to open file "+srcFilename)

	_, err = io.Copy(destFile, srcFile)
	panicOnError(err,
		fmt.Sprintf("Failed to copy data from %s to %s", srcFile.Name(), destFile.Name()))

	err = destFile.Close()
	panicOnError(err, "Failed to close file "+destFile.Name())

	err = srcFile.Close()
	panicOnError(err, "Failed to close file "+srcFile.Name())
}

func mustChmod(filename string, mode os.FileMode) {
	err := os.Chmod(filename, mode)
	panicOnError(err, fmt.Sprintf("Failed to chmod %d %q", mode, filename))
}

// copyDir copies a directory tree over to a new directory.
// Also, dot files and dot directories are skipped.
func mustCopyDir(destDir, srcDir string) error {
	var fullSrcDir string
	// Handle symlinked directories.
	f, err := os.Lstat(srcDir)
	if err == nil && f.Mode()&os.ModeSymlink == os.ModeSymlink {
		fullSrcDir, err = os.Readlink(srcDir)
		if err != nil {
			panic(err)
		}
	} else {
		fullSrcDir = srcDir
	}

	return filepath.Walk(fullSrcDir, func(srcPath string, info os.FileInfo, err error) error {
		// Get the relative path from the source base, and the corresponding path in
		// the dest directory.
		relSrcPath := strings.TrimLeft(srcPath[len(fullSrcDir):], string(os.PathSeparator))
		destPath := path.Join(destDir, relSrcPath)

		// Create a subdirectory if necessary.
		if info.IsDir() {
			err := os.MkdirAll(path.Join(destDir, relSrcPath), 0777)
			if !os.IsExist(err) {
				panicOnError(err, "Failed to create directory")
			}
			return nil
		}

		// Else, just copy it over.
		mustCopyFile(destPath, srcPath)
		return nil
	})
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// empty returns true if the given directory is empty.
// the directory must exist.
func empty(dirname string) bool {
	dir, err := os.Open(dirname)
	if err != nil {
		errorf("error opening directory: %s", err)
	}
	defer dir.Close()
	results, _ := dir.Readdir(1)
	return len(results) == 0
}
