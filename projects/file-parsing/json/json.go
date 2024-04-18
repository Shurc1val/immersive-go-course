package json

import (
	"encoding/json"
	"io"
	"os"
	"parsing/player"
)

type players []player.Player


func (p *players) UnmarshalJSON(data []byte) error {
	var temp []map[string]interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	for _, playerData := range temp {
		player := player.Player{
			Name: playerData["name"].(string),
			HighScore: int(playerData["high_score"].(float64)),
		}
		*p = append(*p, player)
	}
	return nil
}

func decodeJSONData(data []byte) ([]player.Player, error) {
	var p players
	err := json.Unmarshal(data, &p)
	if err != nil {
		return []player.Player(p), err
	}
	
	return []player.Player(p), nil
}

func ReadJSON(filename string) ([]player.Player, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return []player.Player{}, err
	}
	defer jsonFile.Close()

	jsonByteArray, err := io.ReadAll(jsonFile)
	if err != nil {
		return []player.Player{}, err
	}

	return decodeJSONData(jsonByteArray)
}