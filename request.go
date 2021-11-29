package bget

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func Request(url string, headers *map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range *headers {
		req.Header.Add(k, v)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, errors.New("request failed: " + res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// DownloadFileToTemp download url to a temporary file and returns the path to the temporary file
func DownloadFileToTemp(url string, headers *map[string]string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	for k, v := range *headers {
		req.Header.Add(k, v)
	}
	req.Header.Add("Accept", "application/octet-stream")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", errors.New("request failed: " + res.Status)
	}

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
