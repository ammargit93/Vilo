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
		Usage:   "A version control CLI application",
		Version: "0.1.0",
		Commands: []cli.Command{
			{
				Name:    "init",
				Aliases: []string{""},
				Usage:   "Initialises an empty vilo file",
				Action:  func(c *cli.Context) error { return InitCommand() },
			},
			{
				Name:    "add",
				Aliases: []string{""},
				Usage:   "Adds files to Staging area.",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{Name: "files", Usage: "File input paths"},
				},
				Action: func(c *cli.Context) error {
					var filePaths = strings.Split(c.StringSlice("files")[0], ",")
					if len(filePaths) == 0 {
						fmt.Println("No files specified to add")
						return nil
					}
					return AddCommand(filePaths)
				},
			},
			{
				Name:    "commit",
				Aliases: []string{""},
				Usage:   "Commit files .vilo",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "Message", Usage: "commit message"},
				},
				Action: func(c *cli.Context) error {
					commitMsg := c.String("Message")
					if commitMsg == "" {
						fmt.Println("Please provide a commit message using --Message flag")
						return nil
					}

					return CommitCommand(commitMsg)
				},
			},
		},
	}

	log.Fatal(app.Run(os.Args))

}
