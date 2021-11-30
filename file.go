package bget

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func IsQualifiedAsset(name string) bool {
	if IsExecutable(name) {
		return true
	}
	if IsCompressedFile(name) {
		return true
	}
	return false
}

func IsCompressedFile(file string) bool {
	lastExt := filepath.Ext(file)
	firstExt := filepath.Ext(strings.TrimSuffix(file, lastExt))
	ext := firstExt + lastExt
	switch ext {
	case ".tar.gz", ".tgz", ".tar.bz2", ".tbz2", ".tar.xz", ".txz", ".zip":
		return true
	default:
		return false
	}
}

func IsExecutable(file string) bool {
	ext := filepath.Ext(file)

	if runtime.GOOS == "windows" {
		return ext == ".exe"
	}

	return ext == ""
}

type File struct {
	Path string
	Size int64
}

// ReadDir read files recursively in a dir
func ReadDir(dir string) ([]File, error) {
	files := make([]File, 0)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		files = append(files, File{Path: path, Size: info.Size()})
		return nil
	})
	return files, err
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func HumanSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d Bytes", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%d KB", size/1024)
	} else {
		return fmt.Sprintf("%.2f MB", float64(size)/1024/1024)
	}
}
