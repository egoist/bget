package bget

import (
	"strings"
	"testing"
)

func TestGetLatestRelease(t *testing.T) {
	gh := NewGitHub("egoist/doko", "")

	release, err := gh.FetchRelease()
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(release.TagName, "v") {
		t.Fatal("tag name should start with 'v'")
	}
}
