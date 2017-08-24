package cap

import (
	"encoding/binary" 
	"crypto/rand"  
	"io"  
)

func NewId() int64{
	b := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {  
		return 0 
	}  
	result := binary.LittleEndian.Uint64(b)
	return int64(result)
}