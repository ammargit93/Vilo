package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	// "crypto/sha256"
	// "encoding/hex"

	"github.com/urfave/cli"
)

type Node struct {
	hash      string
	data      interface{}
	commitMsg string
	next      *Node
	currTime  time.Time
	nodeName  string
}

func NewNode(hash string, data interface{}, commitMsg string, next *Node, currTime time.Time, nodeName string) Node {
	return Node{
		hash:      hash,
		data:      data,
		commitMsg: commitMsg,
		next:      next,
		currTime:  currTime,
		nodeName:  nodeName,
	}
}

func AddNode(head *Node, newnode Node) *Node {
	h := head
	for {
		if head.next == nil {
			head.next = &newnode
			break
		}
		head = head.next
		fmt.Println("loop")
	}
	head = h
	return head
}

var (
	StagingArea = []string{}
)

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
					headPath := ".vilo/HEAD"
					StagingAreaPath := ".vilo/stage.json"
					if _, err := os.Stat(headPath); err != nil {
						if os.IsNotExist(err) {
							_, err = os.Create(headPath)
							if err != nil {
								fmt.Println(err)
							}
						} else {
							fmt.Println("HEAD already exists")
						}
					}
					if _, err := os.Stat(StagingAreaPath); err != nil {
						if os.IsNotExist(err) {
							_, err = os.Create(StagingAreaPath)
							if err != nil {
								fmt.Println(err)
							}
						} else {
							fmt.Println("stage.json already exists")
						}
					}
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
					fmt.Println(StagingArea)
					fmt.Println("Commit message:", commitMsg)
					for _, f := range StagingArea {
						file, err := os.Open(f)
						if err != nil {
							log.Fatalf("Error opening file %s: %v\n", f, err)
							return err
						}
						defer file.Close()
						hash := sha256.New()
						if _, err := io.Copy(hash, file); err != nil {
							log.Fatalf("Error hashing file %s: %v\n", f, err)
							return err
						}
						fmt.Printf("SHA-256 checksum of %s: %x\n", f, hash.Sum(nil))
					}
					fmt.Println("Commit successful!")
					f, err := os.OpenFile(".vilo/stage.json", os.O_WRONLY|os.O_TRUNC, 0644)
					if err != nil {
						fmt.Println("Error deleting content of the File")
					}
					f.Close()
					return nil
				},
			},
		},
	}

	log.Fatal(app.Run(os.Args))

}
