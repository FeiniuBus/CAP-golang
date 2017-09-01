package cap

import (
	"crypto/rand"
	"encoding/binary"
	"io"
)

// NewID Generate an int64 unique id.
func NewID() int64 {
	b := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return 0
	}
	result := binary.LittleEndian.Uint64(b)
	returnValue := int64(result)
	if returnValue < 0 {
		returnValue = returnValue * -1
	}
	return returnValue
}
