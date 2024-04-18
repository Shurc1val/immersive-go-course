package customBinary

import (
	"encoding/binary"
	"fmt"
	"math"
	"parsing/player"
	"testing"

	"github.com/stretchr/testify/assert"
)

func signedIntToByte(num int32, bigEndian bool) []byte {
	numByte := make([]byte, 4)
	if int(num) > twosComplimentZero {
			num = int32(int(num) - 2*(twosComplimentZero))
		}
	if bigEndian {
		binary.BigEndian.PutUint32(numByte, uint32(math.Abs(float64(num))))
	} else {
		binary.LittleEndian.PutUint32(numByte, uint32(math.Abs(float64(num))))
	}
	return numByte
}

func TestDecodeCustomBinary(t *testing.T) {
	t.Run("Big endian", func(t *testing.T) {
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
				HighScore: 1,
			},
		}
		test_data := []byte{254, 255}
		for _, player := range expected {
			test_data = append(test_data, signedIntToByte(int32(player.HighScore), true)...)
			test_data = append(test_data, []byte(player.Name)...)
			test_data = append(test_data, byte(0))
		}
		fmt.Println(test_data)

		output, err := decodeCustomBinary(test_data)
		if err != nil {
			t.Error("unexpected error")
		}

		assert.Equal(t, expected, output)
	})
	t.Run("Little endian", func(t *testing.T) {
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
				HighScore: 1,
			},
		}
		test_data := []byte{255, 254}
		for _, player := range expected {
			test_data = append(test_data, signedIntToByte(int32(player.HighScore), false)...)
			test_data = append(test_data, []byte(player.Name)...)
			test_data = append(test_data, byte(0))
		}
		fmt.Println(test_data)

		output, err := decodeCustomBinary(test_data)
		if err != nil {
			t.Error("unexpected error")
		}

		assert.Equal(t, expected, output)
	})
}