package bget

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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
		out, err := exec.Command("unzip", "-d", tempDir, file).CombinedOutput()

		if err != nil {
			return "", errors.New(string(out))
		}
	} else {
		out, err := exec.Command("tar", "-xvzf", file, "-C", tempDir).CombinedOutput()

		if err != nil {
			return "", errors.New(string(out))
		}
	}

	// println("extracted to", tempDir)
	return tempDir, nil
}
