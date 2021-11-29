package main

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/egoist/bget"
)

func main() {
	args, err := ParseArgs(os.Args[1:])

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = handle(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handle(args *AppArgs) error {
	gh := bget.NewGitHub(args.Repo, args.GitHubToken)

	stopSpinner := bget.ShowSpinnerWhile("Fetching releases")
	release, err := gh.FetchRelease()
	stopSpinner()
	if err != nil {
		return err
	}

	var assetsForPrompt []string
	assets := make(map[int]bget.GithubReleaseAsset)
	for index, asset := range release.Assets {
		if !bget.IsQualifiedAsset(asset.Name) {
			continue
		}
		assetsForPrompt = append(assetsForPrompt, asset.Name+" ("+bget.HumanSize(asset.Size)+")")
		assets[index] = asset
	}

	if len(assetsForPrompt) == 0 {
		return fmt.Errorf("no releases in this repo")
	}

	var assetIndex int

	err = survey.AskOne(&survey.Select{
		Message: "Select an asset",
		Options: assetsForPrompt,
	}, &assetIndex)

	if err != nil {
		return err
	}

	if assetIndex < 0 || assetIndex >= len(assets) {
		return fmt.Errorf("invalid asset index")
	}

	asset := assets[assetIndex]

	err = gh.DownloadAndInstallAsset(&asset, &bget.InstallOpts{
		BinDir:  args.BinDir,
		BinName: args.BinName,
	})
	return err
}
