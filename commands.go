package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

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

	hash := sha256.Sum256([]byte(commitMsg))
	f, _ = os.OpenFile(".vilo/HEAD", os.O_WRONLY|os.O_TRUNC, 0644)
	f.Close()
	f, _ = os.OpenFile(".vilo/HEAD", os.O_WRONLY|os.O_APPEND, 0644)
	f.WriteString(hex.EncodeToString(hash[:]))
	f.Close()
	commitDir := ".vilo/objects/" + hex.EncodeToString(hash[:]) + "/"
	os.MkdirAll(commitDir, 0755)

	for _, file := range StagingArea {
		fileName := filepath.Base(file)
		outputPath := commitDir + fileName + ".enc"
		err := EncryptAndCompress(file, outputPath, key)
		if err != nil {
			fmt.Printf("Error encrypting file %s: %v\n", file, err)
			continue
		}
		CreateFile(outputPath)
	}

	fmt.Println("Commit successful!")
	DeleteJSONContent()
	return nil
}

func PushCommand(link string) error {
	headFile, _ := os.Open(".vilo/HEAD")
	defer headFile.Close()
	hashBytes, _ := io.ReadAll(headFile)
	hashString := strings.TrimSpace(string(hashBytes))
	commitDir := filepath.Join(".vilo", "objects", hashString)
	err := filepath.Walk(commitDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error accessing path:", path, err)
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".enc") {
			fmt.Println("Pushing file:", path)
			file, err := os.Open(path)
			if err != nil {
				fmt.Println("Error opening file:", err)
				return err
			}
			defer file.Close()

		}
		return nil
	})

	fmt.Println("Push completed successfully!")
	return err
}
