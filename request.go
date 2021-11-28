package bget

import (
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func Request(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// DownloadFileToTemp download url to a temporary file and returns the path to the temporary file
func DownloadFileToTemp(url string) (string, error) {
	// Download the file from `url` and save it to a temp file
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	ext := filepath.Ext(url)
	// Create a temp file
	tmpFile, err := ioutil.TempFile("", "*"+ext)
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	// Write the body to file
	_, err = io.Copy(tmpFile, res.Body)
	if err != nil {
		return "", err
	}

	filepath := tmpFile.Name()
	// fmt.Printf("Downloaded file to temp %s\n", filepath)

	return filepath, nil
}
