package main

import (
	"fmt"
	"os"
	"parsing/csv"
	customBinary "parsing/custom_binary"
	"parsing/json"
	"parsing/player"
	repeatedJson "parsing/repeated_json"
	"sort"
)

const (
	jsonFile string = "examples/json.txt"
	repeatedJsonFile string = "examples/repeated-json.txt"
	csvFile string = "examples/data.csv"
)
var customBinaryFiles = []string{"examples/custom-binary-be.bin", "examples/custom-binary-le.bin"}

func sortPlayersByHighScore(players []player.Player) {
	sort.Slice(players, func(i, j int) bool {
		return players[i].HighScore >= players[j].HighScore
	})
}

func main() {
	players, err := json.ReadJSON(jsonFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("\nJSON")
	sortPlayersByHighScore(players)
	fmt.Println("  High score:", players[0])
	fmt.Println("  Low score:", players[len(players) - 1])

	players, err = repeatedJson.ReadJSON(repeatedJsonFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("\nRepeated JSON")
	sortPlayersByHighScore(players)
	fmt.Println("  High score:", players[0])
	fmt.Println("  Low score:", players[len(players) - 1])

	players, err = csv.ReadCSV(csvFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("\nCSV")
	sortPlayersByHighScore(players)
	fmt.Println("  High score:", players[0])
	fmt.Println("  Low score:", players[len(players) - 1])

	for _, file := range customBinaryFiles {
		players, err = customBinary.ReadCustomBinary(file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("\nCustom binary")
		sortPlayersByHighScore(players)
		fmt.Println("  High score:", players[0])
		fmt.Println("  Low score:", players[len(players) - 1])
	}

}