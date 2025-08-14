package main

import (
	"fmt"
	"os"
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
