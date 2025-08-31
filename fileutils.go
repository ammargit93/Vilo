package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

func DeleteJSONContent() {
	f, err := os.OpenFile(".vilo/stage.json", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error deleting content of the File")
	}
	f.Close()
}

func CreateFile(file string) {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			_, err = os.Create(file)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println(file + "already exists")
		}
	}
}

func SafeSplit(path string) string {
	path = filepath.ToSlash(path)
	parts := strings.Split(path, "/")
	if len(parts) > 3 {
		return strings.Join(parts[3:], "/")
	}
	return path
}

func IsHiddenFile(filename string) (bool, error) {

	if runtime.GOOS == "windows" {
		pointer, err := syscall.UTF16PtrFromString(filename)
		if err != nil {
			return false, err
		}
		attributes, err := syscall.GetFileAttributes(pointer)
		if err != nil {
			return false, err
		}
		return attributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0, nil
	} else {
		// unix/linux file or directory that starts with . is hidden
		if filename[0:1] == "." {
			return true, nil
		}
	}
	return false, nil
}

func ScanRecursively(blobName, commitDir string) error {
	cwd, _ := os.Getwd()
	return filepath.Walk(blobName, func(path string, info fs.FileInfo, err error) error {

		b, ok := IsHiddenFile(path)
		if ok != nil {
			return ok
		}
		if info.IsDir() || b {
			return nil
		}
		relPath, _ := filepath.Rel(cwd, path)
		actual := filepath.Join(commitDir, relPath+".enc")

		if err = os.MkdirAll(filepath.Dir(actual), 0755); err == nil {
			EncryptAndCompress(path, actual, key)
		}
		return err
	})
}
