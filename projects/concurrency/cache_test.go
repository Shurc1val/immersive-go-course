package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	c := NewCache[int, string](4)
	c.Put(1, "Ham")
	c.Put(2, "Cheese")
	c.Put(3, "Onion")
	c.Put(7, "Gherkin")

	t.Run("Get existing element", func(t *testing.T) {
		expected := "Ham"
		output, ok := c.Get(1)

		assert.Equal(t, expected, *output)
		assert.Equal(t, ok, true)
	})
	t.Run("Get non-existent element", func(t *testing.T) {
		output, ok := c.Get(0)

		if output != nil {
			t.Errorf("Output was not nil as expected: %v", output)
		}
		assert.Equal(t, ok, false)
	})
}

func TestPut(t *testing.T) {
	c := NewCache[int, string](4)

	elems := []struct {
		key int
		val string
	}{{4, "Eggs"}, {7, "Green"}, {111, "Albuquerque"}}

	for _, elem := range elems {
		t.Run("Put non-existent elements", func(t *testing.T) {
			c.Put(elem.key, elem.val)
			output, ok := c.Get(elem.key)
			assert.Equal(t, elem.val, *output)
			assert.Equal(t, ok, true)
		})
	}

	for _, elem := range elems {
		t.Run("Put existent elements", func(t *testing.T) {
			c.Put(elem.key, "FUDGE")
			output, ok := c.Get(elem.key)
			assert.Equal(t, "FUDGE", *output)
			assert.Equal(t, ok, true)
		})
	}
}

func TestRetrieveStats(t *testing.T) {
	c := NewCache[int, string](4)
	c.Put(1, "Ham")
	c.Put(2, "Cheese")
	c.Put(3, "Onion")
	c.Put(7, "Gherkin")
	c.Get(1)
	c.Get(3)
	c.Put(1, "Salad")
	c.Get(20)
	c.Get(30)
	c.Get(1)
	c.Put(5, "Lemons")
	c.Put(6, "Apple")
	fmt.Println(c.lru_tracker)

	expected := []any{0.6, 4, 12}
	output := make([]any, 3)
	output[0], output[1], output[2] = c.RetrieveStats()
	for i, val := range output {
		if val != expected[i] {
			t.Errorf("Output, %v, did not match expected, %v", val, expected[i])
		}
	}
}
