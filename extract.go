package bget

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Extract extracts the file to a temporary directory
func Extract(file string) (string, error) {
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}
	os.MkdirAll(tempDir, 0755)

	// extract
	ext := filepath.Ext(file)
	if ext == ".zip" {
		bytes, err := exec.Command(wrapCommand("unzip"), "-d", tempDir, file).CombinedOutput()
		if err != nil {
			output := string(bytes)
			if output != "" {
				return "", errors.New(output)
			}
			return "", err
		}
	} else {
		bytes, err := exec.Command(wrapCommand("tar"), "-xvzf", file, "-C", tempDir).CombinedOutput()
		if err != nil {
			output := string(bytes)
			if output != "" {
				return "", errors.New(output)
			}
			return "", err
		}
	}

	// println("extracted to", tempDir)
	return tempDir, nil
}

func wrapCommand(command string) string {
	if runtime.GOOS == "windows" {
		return command + ".exe"
	}
	return command
}
