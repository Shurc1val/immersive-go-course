package cmd

import (
	"fmt"
	"io/fs"
	"os"
)

func getDirectoryContents(directory string) ([]string, error) {
	var contents []string
	sub_dirs, err := os.ReadDir(directory)
	if err != nil {
		return contents, err
	}
	for _, sub_dir := range sub_dirs {
		contents = append(contents, sub_dir.Name())
	}
	return contents, nil
}

func Execute() {
	args := os.Args[1:]
	var directory string
	if len(args) > 0 {
		directory = args[0]
	} else {
		var err error
		directory, err = os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	dirContents, err := getDirectoryContents(directory)
	switch err.(type) {
	case *fs.PathError:
		potentialFile := directory
		_, err := os.Stat(potentialFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(potentialFile)
	case nil:
		for _, name := range dirContents {
			fmt.Println(name)
		}
	default:
		fmt.Println(err)
		os.Exit(1)
	}

}