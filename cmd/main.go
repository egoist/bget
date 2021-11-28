package main

import (
	"fmt"
	"os"
	"path/filepath"

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
	gh := bget.GitHub{}

	gh.ParseRepo(args.Repo)

	stopSpinner := bget.ShowSpinnerWhile("Fetching releases")
	release, err := gh.FetchRelease()
	if err != nil {
		stopSpinner()
		return err
	}
	stopSpinner()

	var assetsForPrompt []string
	assets := make(map[string]bget.GithubReleaseAsset)
	for _, asset := range release.Assets {
		if !bget.IsQualifiedAsset(asset.Name) {
			continue
		}
		assetsForPrompt = append(assetsForPrompt, asset.Name)
		assets[asset.Name] = asset
	}

	assetName := ""

	err = survey.AskOne(&survey.Select{
		Message: "Select an asset",
		Options: assetsForPrompt,
	}, &assetName)

	if err != nil {
		return err
	}

	if assetName == "" {
		return fmt.Errorf("no asset selected")
	}

	asset := assets[assetName]

	stopSpinner = bget.ShowSpinnerWhile("Downloading " + assetName)
	tempFile, err := bget.DownloadFileToTemp(asset.DwnloadURL)
	if err != nil {
		stopSpinner()
		return err
	}
	stopSpinner()

	binFile := ""

	if bget.IsCompressedFile(assetName) {
		tempDir, err := bget.Extract(tempFile)
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
	if args.BinDir != "" {
		binDir = args.BinDir
	}
	if args.BinName != "" {
		binName = args.BinName
	}
	err = InstallBin(binFile, binName, binDir)
	if err != nil {
		return err
	}
	return nil
}

func InstallBin(src string, binName string, binDir string) error {
	dest := filepath.Join(binDir, binName)
	err := os.Rename(src, dest)
	if err != nil {
		return err
	}

	// Make dest executable
	err = os.Chmod(dest, 0755)

	println("Installed to:", dest)

	return err
}

func GetBinFromDir(dir string) (string, error) {
	files, err := bget.ReadDir(dir)
	if err != nil {
		return "", err
	}

	var largest int64 = 0
	filepath := ""

	for _, file := range files {

		if !bget.IsExecutable(file.Path) {
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
