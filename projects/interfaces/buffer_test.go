package buffer

import (
	"testing"
)

func TestBytes(t *testing.T) {
	test_byte_slice := []byte("Testing this thing.")
	b := NewBuffer(test_byte_slice)
	for i, val := range b.Bytes() {
		if val != test_byte_slice[i] {
			t.Fatal("b.Bytes did not return the bytes the buffer was initialised with")
		}
	}
}

func TestWrite(t *testing.T) {
	test_byte_slices := [][]byte{
		[]byte("Testing this thing."),
		[]byte("STILL testing this thing!!"),
	}
	total_byte_slice := append(test_byte_slices[0], test_byte_slices[1]...)

	b := NewBuffer(test_byte_slices[0])
	b.Write(test_byte_slices[1])
	for i, val := range b.Bytes() {
		if val != total_byte_slice[i] {
			t.Fatal("b.Write did not work as intended.")
		}
	}
}

func TestRead(t *testing.T) {
	t.Run("Slice larger than bytes contained", func (t *testing.T) {
		test_byte_slice := []byte("If you like pina coladas")
		p := make([]byte, len(test_byte_slice) + 1)

		b := NewBuffer(test_byte_slice)
		b.Read(p)
		for i, val := range test_byte_slice {
			if val != p[i] {
				t.Fatal("Not all bytes were read when there was space available.")
			}
		}
	})

	t.Run("Slice smaller than bytes contained", func (t *testing.T) {
		test_byte_slice := []byte("If you like pina coladas")
		p := make([]byte, len(test_byte_slice) - 2)

		b := NewBuffer(test_byte_slice)
		b.Read(p)
		for i, val := range p {
			if val != test_byte_slice[i] {
				t.Fatal("Bytes read incorrectly.")
			}
		}
		b.Read(p)
		for i, val := range p[:len(test_byte_slice) - len(p)] {
			if val != test_byte_slice[len(p) + i] {
				t.Fatal("Remaining bytes read incorrectly.")
			}
		}
	})

}

type testWriter [][]byte

func (tw testWriter) Write(p []byte) (int, error) {
	for _, val := range p {
		tw[0] = append(tw[0], val)
	}
	return 0, nil
}

func TestFilteringPipe(t *testing.T) {
	tests := map[string]struct {
		input []byte
		expectedOutput []byte
	}{
		"no digits" : {
			input: []byte("This is just text."),
			expectedOutput: []byte("This is just text."),
		},
		"some digits" : {
			input: []byte("It was hard 4 me 2 talk 2u."),
			expectedOutput: []byte("It was hard  me  talk u."),
		},
		"all digits" : {
			input: []byte("45736"),
			expectedOutput: []byte{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tw := testWriter{[]byte{}}
			testFP := NewFilteringPipe(tw)
			testFP.Write(test.input)
			for i, val := range tw[0] {
				if val != test.expectedOutput[i] {
					t.Fatalf("Writer was called with %q, instead of %q as expected, for input of %q", tw[0], test.expectedOutput, test.input)
				}
			}
		})
	}
}