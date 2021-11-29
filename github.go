package bget

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

func NewGitHub(input string, token string) *GitHub {
	gh := &GitHub{}

	gh.New(input, token)

	return gh
}

type GitHub struct {
	Owner   string
	Repo    string
	TagName string
	Token   string
}

type GitHubRelease struct {
	TagName string               `json:"tag_name"`
	Assets  []GithubReleaseAsset `json:"assets"`
}

type GithubReleaseAsset struct {
	DwnloadURL string `json:"url"`
	Size       int64  `json:"size"`
	Name       string `json:"name"`
}

func (gh *GitHub) New(input string, token string) {
	repoAndTag := strings.Split(input, "#")
	ownAndRepo := strings.Split(repoAndTag[0], "/")
	gh.Owner = ownAndRepo[0]
	gh.Repo = ownAndRepo[1]
	if len(repoAndTag) == 2 {
		gh.TagName = repoAndTag[1]
	}
	gh.Token = os.Getenv("GITHUB_TOKEN")
	if token != "" {
		gh.Token = token
	}
}

func (gh *GitHub) GetHeaders() *map[string]string {
	headers := map[string]string{
		"User-Agent": "bget",
	}
	if gh.Token != "" {
		headers["Authorization"] = fmt.Sprintf("token %s", gh.Token)
	}
	return &headers
}

func (gh *GitHub) FetchRelease() (GitHubRelease, error) {
	url := ""
	if gh.TagName == "" {
		url = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", gh.Owner, gh.Repo)
	} else {
		url = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", gh.Owner, gh.Repo, gh.TagName)
	}

	data, err := Request(url, gh.GetHeaders())
	if err != nil {
		return GitHubRelease{}, err
	}

	var release GitHubRelease
	err = json.Unmarshal(data, &release)

	if err != nil {
		return GitHubRelease{}, err
	}

	return release, nil
}

type InstallOpts struct {
	BinDir  string
	BinName string
}

func (gh *GitHub) DownloadAndInstallAsset(asset *GithubReleaseAsset, installOpts *InstallOpts) error {
	stopSpinner := ShowSpinnerWhile("Downloading " + asset.Name)
	tempFile, err := DownloadFileToTemp(asset.DwnloadURL, gh.GetHeaders())
	stopSpinner()

	if err != nil {
		return err
	}

	binFile := ""

	if IsCompressedFile(asset.Name) {
		tempDir, err := Extract(tempFile)
		if err != nil {
			return err
		}
		binFile, err = GetBinFromDir(tempDir)
		if err != nil {
			return err
		}

	} else {
		binFile = tempFile
	}

	binDir := "/usr/local/bin"
	binName := gh.Repo
	if installOpts.BinDir != "" {
		binDir = installOpts.BinDir
	}
	if installOpts.BinName != "" {
		binName = installOpts.BinName
	}
	err = InstallBin(binFile, binName, binDir)
	if err != nil {
		return err
	}
	return nil

}

func InstallBin(src string, binName string, binDir string) error {
	dest := filepath.Join(binDir, binName)

	if PathExists(dest) {
		overwrite := false
		err := survey.AskOne(&survey.Confirm{
			Message: "Bin already exists. Overwrite?",
		}, &overwrite)
		if err != nil {
			return err
		}
		if !overwrite {
			return fmt.Errorf("aborted")
		}
	}

	os.MkdirAll(binDir, 0755)

	// Try to move it
	err := os.Rename(src, dest)

	// TODO: handle permission error
	// Sometimes you need sudo access to the binDir
	if err != nil {
		return err
	}

	// Make dest executable
	err = os.Chmod(dest, 0755)

	println("Installed to:", dest)

	return err
}

func GetBinFromDir(dir string) (string, error) {
	files, err := ReadDir(dir)
	if err != nil {
		return "", err
	}

	var largest int64 = 0
	filepath := ""

	for _, file := range files {

		if !IsExecutable(file.Path) {
			continue
		}

		if file.Size > largest {
			largest = file.Size
			filepath = file.Path
		}
	}

	if filepath == "" {
		return "", fmt.Errorf("no executable file found in %s", dir)
	}

	return filepath, nil
}
