package main

import (
	"fmt"
	"log"
	"os"
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
	StagingArea = []os.File{}
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
					return err
				},
			},
		},
	}

	log.Fatal(app.Run(os.Args))

}
