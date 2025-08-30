package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Node struct {
	currHash string
	prevHash *Node
	nextHash *Node
}

var InitialNode Node
var head *Node = nil

func InitCommand() error {

	err := os.MkdirAll(".vilo", 0755)
	if err != nil {
		fmt.Println("Error while creating a .vilo dir, ", err)
	}
	err = os.MkdirAll(".vilo/objects", 0755)
	if err != nil {
		fmt.Println("Error while creating a objects dir, ", err)
	}

	CreateFile(".vilo/HEAD")
	CreateFile(".vilo/history")
	CreateFile(".vilo/stage.json")
	return err
}

func AddCommand(filePaths []string) error {
	DeleteJSONContent()

	for _, file := range filePaths {
		absPath, _ := filepath.Abs(file)
		absPath = strings.TrimSpace(absPath)
		if _, err := os.Stat(absPath); err == nil {
			StagingArea = append(StagingArea, absPath)
			fmt.Println(absPath, "Staged for commit")
		} else {
			fmt.Printf("File does not exist in the staging area: %s", absPath)
		}
	}
	f, err := os.OpenFile(".vilo/stage.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("stage.json doesn't exist, use vilo init first.")
		return err
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "	")
	if err := encoder.Encode(StagingArea); err != nil {
		log.Fatalf("Error encoding data to JSON: %v", err)
	}
	fmt.Println(StagingArea)
	return nil
}

func CommitCommand(commitMsg string) error {

	f, _ := os.Open(".vilo/stage.json")
	decoder := json.NewDecoder(f)
	decoder.Decode(&StagingArea)
	f.Close()

	hashedCommit := GenerateCommitHash(commitMsg, StagingArea)

	if err := os.WriteFile(".vilo/HEAD", []byte(hashedCommit), 0644); err != nil {
		return err
	}
	f, err := os.OpenFile(".vilo/history", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(hashedCommit + " " + commitMsg + "\n"); err != nil {
		return err
	}

	commitDir := ".vilo/objects/" + hashedCommit + "/"
	os.MkdirAll(commitDir, 0755)

	for _, file := range StagingArea {
		fileName := filepath.Base(file)
		fmt.Println(fileName)
		ScanDirRecursively(file, commitDir)

	}

	fmt.Println("Commit successful!")

	DeleteJSONContent()
	return nil
}

func ShowCommits() {
	f, _ := os.ReadFile(".vilo/history")
	fmt.Println(string(f))
}
func safeSplit(path string) string {
	path = filepath.ToSlash(path)
	parts := strings.Split(path, "/")
	if len(parts) > 3 {
		return strings.Join(parts[3:], "/")
	}
	return path
}

func RollBack(commitHash string, backupFolderName string) {
	commitDir := ".vilo/objects/" + commitHash + "/"
	_ = filepath.WalkDir(commitDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		actualPath := safeSplit(path)
		actualPath = strings.TrimSuffix(actualPath, ".enc")

		outputPath := filepath.Join(backupFolderName, actualPath)

		os.MkdirAll(filepath.Dir(outputPath), 0755)

		fmt.Println(actualPath)
		err = DecryptAndDecompress(path, outputPath, key)
		if err != nil {
			fmt.Println("Error restoring", path, ":", err)
		} else {
			fmt.Println("Restored", outputPath)
		}
		return nil
	})

}
