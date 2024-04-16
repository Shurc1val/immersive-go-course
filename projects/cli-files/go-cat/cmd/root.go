package cmd

import (
	"fmt"
	"os"
)

func Execute() {
	args := os.Args[1:]
	var directory string
	if len(args) > 0 {
		directory = args[0]
	} else {
		fmt.Println("No file given.")
		os.Exit(1)
	}
	contents, err := os.ReadFile(directory)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Stdout.Write(contents)
}