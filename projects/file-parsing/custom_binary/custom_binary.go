package customBinary

import (
	"encoding/binary"
	"errors"
	"os"
	"parsing/player"
)

const twosComplimentZero int = 2147483648

func decodeCustomBinary(data []byte) ([]player.Player, error) {
	players := []player.Player{}
	endian := "big"
	if (len(data) > 2) && data[0] == 255 && data[1] == 254 {
		endian = "little"
	} else if (len(data) < 2) || !(data[0] == 254 && data[1] == 255) {
		return players, errors.New("unsupported binary encoding")
	}
	for i := 6; i < len(data); {
		var highScore int
		scoreBytes := data[i-4:i]
		if endian == "big" {
			highScore = int(binary.BigEndian.Uint32(scoreBytes))
		} else {
			highScore = int(binary.LittleEndian.Uint32(scoreBytes))
		}
		if highScore > twosComplimentZero {
			highScore -= 2*twosComplimentZero
		}
		
		nameBytes := []byte{}
		for data[i] != 0 { 
			nameBytes = append(nameBytes, data[i])
			i++
		}
		name := string(nameBytes[:])
		players = append(players, player.Player{
			Name: name,
			HighScore: highScore,
		})
		i += 5
	}
	return players, nil
}

func ReadCustomBinary(filename string) ([]player.Player, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return []player.Player{}, err
	}
	return decodeCustomBinary(data)
}