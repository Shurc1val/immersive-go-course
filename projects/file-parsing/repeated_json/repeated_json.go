package repeatedJson

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"parsing/player"
) 

func decodeJSONData(data []byte) ([]player.Player, error) {
	var players []player.Player
	for _, line := range bytes.Split(data, []byte("\n")) {
		if (len(line) > 0) && (line[0] != []byte("#")[0]) {
			var p player.Player
			err := json.Unmarshal(line, &p)
			if err != nil {
				return players, err
			}
			players = append(players, p)
		}
	}
	return players, nil
}

func ReadJSON(filename string) ([]player.Player, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return []player.Player{}, err
	}
	defer jsonFile.Close()

	fileByteArray, err := io.ReadAll(jsonFile)
	if err != nil {
		return []player.Player{}, err
	}
	
	return decodeJSONData(fileByteArray)
}