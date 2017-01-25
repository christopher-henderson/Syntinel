package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// ProjectRoot returns the absolutely path the directory the binary resides in.
func ProjectRoot() string {
	return projectRoot
}

// ScriptDirectory returns the absolute path to the assets/scripts directory.
func ScriptDirectory() string {
	return scripts
}

// BuildDirectory returns the absolute path to the assets/build directory.
func BuildDirectory() string {
	return build
}

// DockerFileDirectory returns the absolute path to the assets/dockerfiles directory.
func DockerFileDirectory() string {
	return dockerfiles
}

func DatabaseDirectory() string {
	return database
}

var projectRoot = func() (path string) {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		// If we errored here, then we deserve to panic.
		panic(err)
	}
	return path
}()
var scripts = func() string {
	return fmt.Sprintf("%v%v%v%v%v%v", projectRoot, string(os.PathSeparator), "assets", string(os.PathSeparator), "scripts", string(os.PathSeparator))
}()
var build = func() string {
	return fmt.Sprintf("%v%v%v%v%v%v", projectRoot, string(os.PathSeparator), "assets", string(os.PathSeparator), "build", string(os.PathSeparator))
}()
var dockerfiles = func() string {
	return fmt.Sprintf("%v%v%v%v%v%v", projectRoot, string(os.PathSeparator), "assets", string(os.PathSeparator), "dockers", string(os.PathSeparator))
}()
var database = func() string {
	return fmt.Sprintf("%v%v%v%v%v%v", projectRoot, string(os.PathSeparator), "assets", string(os.PathSeparator), "database", string(os.PathSeparator))
}()
