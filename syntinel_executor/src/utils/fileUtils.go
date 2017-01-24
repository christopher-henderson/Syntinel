package utils

import (
	"io"
	"os"
)

// FileCopy does a complete copy from the source to the destination on the
// filesystem.
//
// If the source does not exists then this is a no-op.
func FileCopy(source string, destination string) error {
	// Does the source exist?
	if _, err := os.Stat(source); err != nil {
		return err
	}
	// Open the source file.
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()
	// Open the destination file.
	dst, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer dst.Close()
	// Copy source to destination.
	if _, err := io.Copy(dst, src); err != nil {
		return err
	}
	return nil
}

// FileRemove deletes the given file path from the filesystem. If the file
// does not exists, then this is a no-op
func FileRemove(path string) error {
	if _, err := os.Stat(path); err != nil {
		return nil
	}
	if err := os.Remove(path); err != nil {
		return err
	}
	return nil
}
