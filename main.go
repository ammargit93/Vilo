package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli"
)

var (
	StagingArea = []string{}
	key         = []byte("thisis32bytekeythisis32bytekey!!")
)

func main() {

	var app = &cli.App{
		Name:    "vilo",
		Usage:   "A backup CLI tool",
		Version: "0.1.0",
		Commands: []cli.Command{
			{
				Name:   "init",
				Usage:  "Initialises an empty vilo file",
				Action: func(c *cli.Context) error { return InitCommand() },
			},
			{
				Name:  "add",
				Usage: "Adds files to Staging area.",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{Name: "files", Usage: "File input paths"},
				},
				Action: func(c *cli.Context) error {
					var filePaths = strings.Split(c.StringSlice("files")[0], ",")
					if filePaths[0] == "." {
						pathToIterate, _ := os.Getwd()
						entries, err := os.ReadDir(pathToIterate)
						if err != nil {
							log.Fatal(err)
						}
						for _, e := range entries {
							filePaths = append(filePaths, e.Name())
						}
					}
					if _, err := os.Stat(".viloignore"); err == nil {
						fmt.Println("File exists")
					} else if os.IsNotExist(err) {
						fmt.Println("File does not exist")
					} else {
						fmt.Println("Error checking file:", err)
					}

					if len(filePaths) == 0 {
						fmt.Println("No files specified to add")
						return nil
					}
					return AddCommand(filePaths)
				},
			},
			{
				Name:  "commit",
				Usage: "Commit files .vilo",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "Message", Usage: "commit message"},
					&cli.StringFlag{Name: "cloud", Usage: "upload to cloud"},
				},
				Action: func(c *cli.Context) error {
					commitMsg := c.String("Message")
					cloudProvider := c.String("cloud")
					if commitMsg == "" {
						fmt.Println("Please provide a commit message using --Message flag")
						return nil
					}
					fmt.Println(cloudProvider)

					return CommitCommand(commitMsg)
				},
			},
			{
				Name:  "show",
				Usage: "Displays all commits.",
				Flags: []cli.Flag{},
				Action: func(c *cli.Context) error {
					ShowCommits()
					return nil
				},
			},
			{
				Name:  "rollback",
				Usage: "Rolls back to previous or future commits.",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "hash", Usage: "commit hash"},
					&cli.StringFlag{Name: "foldername", Usage: "backup folder name"},
				},
				Action: func(c *cli.Context) error {
					commitHash := c.String("hash")
					backupFolderName := c.String("foldername")
					if commitHash == "" {
						fmt.Println("Please provide a commit hash using --hash flag")
						return nil
					}
					if backupFolderName == "." {
						backupFolderName, err := os.Getwd()
						if err != nil {
							return nil
						}
						RollBack(commitHash, backupFolderName)
					} else {
						RollBack(commitHash, backupFolderName)
					}

					return nil
				},
			},
			{
				Name:  "status",
				Usage: "Displays staged files.",
				Flags: []cli.Flag{},
				Action: func(c *cli.Context) error {
					Status()
					return nil
				},
			},
			// {
			// 	Name:    "watch-mode",
			// 	Aliases: []string{""},
			// 	Usage:   "Backs up the project, add+commit",
			// 	Flags:   []cli.Flag{},
			// 	Action: func(c *cli.Context) error {
			// 		// BackupCommand()
			// 		return nil
			// 	},
			// },
		},
	}

	log.Fatal(app.Run(os.Args))

}
