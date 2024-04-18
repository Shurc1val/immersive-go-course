package csv

import (
	"encoding/csv"
	"io"
	"os"
	"parsing/player"
	"strconv"
)

func decodeCSV(r io.Reader) ([]player.Player, error) {
	players := []player.Player{}
	csvReader := csv.NewReader(r)
	records, err := csvReader.ReadAll()
	if err != nil {
		return players, err
	}
	for _, record := range records[1:] {
		highScore, err := strconv.Atoi(record[1])
		if err != nil {
			return players, err
		}
		players = append(players, player.Player{
			Name: record[0],
			HighScore: highScore,
		})
	}
	return players, nil
}

func ReadCSV(filename string) ([]player.Player, error) {
	csvFile, err := os.Open(filename)
	if err != nil {
		return []player.Player{}, err
	}
	defer csvFile.Close()

	return decodeCSV(csvFile)
}