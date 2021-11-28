package bget

import (
	"encoding/json"
	"fmt"
	"strings"
)

type GitHub struct {
	Owner   string
	Repo    string
	TagName string
}

type GitHubRelease struct {
	TagName string               `json:"tag_name"`
	Assets  []GithubReleaseAsset `json:"assets"`
}

type GithubReleaseAsset struct {
	DwnloadURL string `json:"browser_download_url"`
	Size       int64  `json:"size"`
	Name       string `json:"name"`
}

func (gh *GitHub) ParseRepo(input string) {
	repoAndTag := strings.Split(input, "#")
	ownAndRepo := strings.Split(repoAndTag[0], "/")
	gh.Owner = ownAndRepo[0]
	gh.Repo = ownAndRepo[1]
	if len(repoAndTag) == 2 {
		gh.TagName = repoAndTag[1]
	}
}

func (gh *GitHub) FetchRelease() (GitHubRelease, error) {
	url := ""
	if gh.TagName == "" {
		url = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", gh.Owner, gh.Repo)
	} else {
		url = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", gh.Owner, gh.Repo, gh.TagName)
	}
	data, err := Request(url)
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
