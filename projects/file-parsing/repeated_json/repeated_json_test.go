package repeatedJson

import (
	"parsing/player"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeJSONData(t *testing.T) {
	data := []byte(`{"name": "Hugh", "high_score": 10}
{"name": "Percival", "high_score": 30}
# This is a comment, and should be ignored
{"name": "Morris", "high_score": -1}
`)
	expected := []player.Player{
		{
			Name: "Hugh",
			HighScore: 10,
		},
		{
			Name: "Percival",
			HighScore: 30,
		},
		{
			Name: "Morris",
			HighScore: -1,
		},
	}

	output, _ := decodeJSONData(data)

	assert.Equal(t, expected, output)
}