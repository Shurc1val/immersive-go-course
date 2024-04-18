package csv

import (
	"bytes"
	"parsing/player"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeCSV(t *testing.T) {
	test_reader := bytes.NewBufferString(`name,high score
Hugh,10
Percival,30
Morris,-1`)
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

	output, _ := decodeCSV(test_reader)
	assert.Equal(t, expected, output)
}