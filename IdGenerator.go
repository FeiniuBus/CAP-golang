package cap

import (
	"os"
	"encoding/binary"  
)

func NewId() int64{
	f, _ := os.OpenFile("/dev/urandom", os.O_RDONLY, 0)
	b := make([]byte, 16)
	f.Read(b)
	f.Close()
	result := binary.LittleEndian.Uint64(b)
	return int64(result)
}