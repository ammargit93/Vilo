package main

import (
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

var (
	StagingArea = []string{}
	key         = []byte("thisis32bytekeythisis32bytekey!!")
)

func EncryptAndCompress(inputPath, outputPath string, key []byte) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return err
	}
	if _, err := outputFile.Write(iv); err != nil {
		return err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	gzipWriter := gzip.NewWriter(cipher.StreamWriter{S: stream, W: outputFile})
	defer gzipWriter.Close()
	_, err = io.Copy(gzipWriter, inputFile)
	return err
}

func DecryptAndDecompress(inputPath, outputPath string, key []byte) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(inputFile, iv); err != nil {
		return err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	stream := cipher.NewCFBDecrypter(block, iv)
	gzipReader, err := gzip.NewReader(cipher.StreamReader{S: stream, R: inputFile})
	if err != nil {
		return err
	}
	defer gzipReader.Close()
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, gzipReader)
	return err
}

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

func main() {

	var app = &cli.App{
		Name:    "vilo",
		Usage:   "A version control CLI application",
		Version: "0.1.0",
		Commands: []cli.Command{
			{
				Name:    "init",
				Aliases: []string{""},
				Usage:   "Initialises a empty vilo file",
				Action: func(c *cli.Context) error {
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
				},
			},
			{
				Name:    "add",
				Aliases: []string{""},
				Usage:   "adds files to Staging area.",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:  "files",
						Usage: "File input paths",
					},
				},
				Action: func(c *cli.Context) error {
					var filePaths = strings.Split(c.StringSlice("files")[0], ",")
					if len(filePaths) == 0 {
						fmt.Println("No files specified to add")
						return nil
					}
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
						fmt.Println("stage.json doesnt exist, use vilo init first.")
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
				},
			},
			{
				Name:    "commit",
				Aliases: []string{""},
				Usage:   "Commit files .vilo",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "Message",
						Usage: "commit message",
					},
				},
				Action: func(c *cli.Context) error {
					commitMsg := c.String("Message")
					if commitMsg == "" {
						fmt.Println("Please provide a commit message using --Message flag")
						return nil
					}
					f, _ := os.Open(".vilo/stage.json")
					decoder := json.NewDecoder(f)
					decoder.Decode(&StagingArea)
					f.Close()

					hash := sha256.Sum256([]byte(commitMsg))
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
				},
			},
		},
	}

	log.Fatal(app.Run(os.Args))

}
