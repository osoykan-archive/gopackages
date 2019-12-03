package main

import (
	"./src"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Please provide a folder name and prefix!")
		os.Exit(-1)
	}

	err := filerenamer.RenameFiles(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Printf("There was and error while processing the command: %s", err.Error())
		os.Exit(-1)
		return
	}

	fmt.Print("Success! Files are renamed.")
	os.Exit(0)
}
