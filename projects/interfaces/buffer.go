package buffer

import (
	"io"
	"os"
	"unicode"
)

type OurByteBuffer struct {
	bytes []byte
}

func NewBuffer(b []byte) *OurByteBuffer {
	return &OurByteBuffer{bytes: b}
}

func (bb OurByteBuffer) Bytes() []byte {
	return bb.bytes
}

func (bb *OurByteBuffer) Write(b []byte) {
	bb.bytes = append(bb.bytes, b...)
}

func (bb *OurByteBuffer) Read(b []byte) {
	for i := range len(b) {
		if len(bb.bytes) == 0 {
			break
		}
		b[i] = bb.bytes[0]
		bb.bytes = bb.bytes[1:]
	}
}

type FilteringPipe struct {
	writer io.Writer
}

func NewFilteringPipe(writer io.Writer) *FilteringPipe {
	return &FilteringPipe{writer: writer}
}

func (f FilteringPipe) Write(b []byte) {
	output := []byte{}
	for _, val := range b {
		if !unicode.IsDigit(rune(val)) {
			output = append(output, val)
		}
	}
	f.writer.Write(output)
}

func main() {
	filteringPipe := NewFilteringPipe(os.Stdout)
	filteringPipe.Write([]byte("start=1, end=10"))
}