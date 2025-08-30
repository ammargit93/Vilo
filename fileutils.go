package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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

func ScanDirRecursively(rootDir, commitDir string) error {
	cwd, _ := os.Getwd()

	return filepath.Walk(rootDir, func(path string, info fs.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if info.IsDir() {
			return nil // skip encrypting dirs themselves
		}

		// get relative path inside project (everything after cwd)
		rel, err := filepath.Rel(cwd, path)
		if err != nil {
			return err
		}

		// add ".enc" and join with commit dir
		outputPath := filepath.Join(commitDir, rel+".enc")

		// make sure subdirectories exist before writing
		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return err
		}

		// encrypt file → write to output path
		if err := EncryptAndCompress(path, outputPath, key); err != nil {
			fmt.Printf("Error encrypting file %s: %v\n", path, err)
			return err
		}

		fmt.Println("Encrypted:", path, "→", outputPath)
		return nil
	})
}
