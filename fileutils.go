package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
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

func remove(arr []string, path string) []string {
	var newArr []string
	for _, val := range arr {
		val, _ := filepath.Abs(val)
		if val != path {
			newArr = append(newArr, val)
		}
	}
	return newArr

}

func FilterIgnoredFiles(fileArr []string) []string {
	projectRoot, _ := os.Getwd()
	fileArr = remove(fileArr, projectRoot)
	f, err := os.Open(".viloignore")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		txt := scanner.Text()
		trimmedPath := strings.TrimSpace(txt)
		trimmedAbsPath, _ := filepath.Abs(trimmedPath)
		fileArr = remove(fileArr, trimmedAbsPath)
	}
	return fileArr
}

func FindLatestCommit() int {
	entries, err := os.ReadDir(".vilo/objects")
	if err != nil {
		fmt.Println(err)
	}
	latestCommit := 0
	for _, entry := range entries {
		entryArr := strings.Split(entry.Name(), "_")
		if entryArr[len(entryArr)-1] == "base" {
			continue
		}
		latestCommitNum, _ := strconv.Atoi(entryArr[len(entryArr)-1])
		latestCommit = max(latestCommitNum, latestCommit)
	}

	return latestCommit
}
