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

	stopSpinner = bget.ShowSpinnerWhile("Downloading " + asset.Name)
	tempFile, err := bget.DownloadFileToTemp(asset.DwnloadURL)
	stopSpinner()
	if err != nil {
		return err
	}

	binFile := ""

	if bget.IsCompressedFile(asset.Name) {
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

	if bget.PathExists(dest) {
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
