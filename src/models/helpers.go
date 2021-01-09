package models

import "encoding/binary"

func Uint64Key(k uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, k)
	return b
}

func KeyUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}
