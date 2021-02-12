package endianness

import (
	"math/bits"
	"unsafe"
)

var (
	isLE               = false
	FromLittleEndian16 = ToLittleEndian16
)

func init() {
	val := uint16(1)
	arr := (*[2]byte)(unsafe.Pointer(&val))
	isLE = arr[0] != 0
}

func ToLittleEndian16(v uint16) uint16 {
	if isLE {
		return v
	}
	return bits.ReverseBytes16(v)
}
