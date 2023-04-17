package utils

import (
	"encoding/binary"
	"encoding/hex"
	"strings"
)

// Bytes1 returns the first byte of the little-endian representation of the supplied value.
func Bytes1(val uint64) []byte {
	return bytesN(val, 1)
}

// Bytes2 returns the first two bytes of the little-endian representation of the supplied value.
func Bytes2(val uint64) []byte {
	return bytesN(val, 2)
}

// Bytes4 returns the first four bytes of the little-endian representation of the supplied value.
func Bytes4(x uint64) []byte {
	return bytesN(x, 4)
}

// Bytes8 returns the first eight bytes of the little-endian representation of the supplied value.
func Bytes8(x uint64) []byte {
	return bytesN(x, 8)
}

// Bytes16 returns the first sixteen bytes of the little-endian representation of the supplied value.
func Bytes16(x uint64) []byte {
	return bytesN(x, 16)
}

// Bytes32 returns the first thirty-two bytes of the little-endian representation of the supplied value.
func Bytes32(x uint64) []byte {
	return bytesN(x, 32)
}

// Bytes64 returns the first thirty-two bytes of the little-endian representation of the supplied value.
func Bytes64(x uint64) []byte {
	return bytesN(x, 64)
}

// bytesN returns the first n bytes of the little-endian representation of the supplied value.
// n must be <= 32
func bytesN(val uint64, n uint) []byte {
	bytes := make([]byte, 64)
	binary.LittleEndian.PutUint64(bytes, val)
	return bytes[:n]
}

// ToBytes8 returns an 8-byte array with the supplied value placed in the low-order indices.
func ToBytes8(val []byte) [8]byte {
	var res [8]byte
	copy(res[:], val)
	return res
}

// ToBytes16 returns a 16-byte array with the supplied value placed in the low-order indices.
func ToBytes16(val []byte) [16]byte {
	var res [16]byte
	copy(res[:], val)
	return res
}

// ToBytes32 returns a 32-byte array with the supplied value placed in the low-order indices.
func ToBytes32(val []byte) [32]byte {
	var res [32]byte
	copy(res[:], val)
	return res
}

// ToBytes48 returns a 48-byte array with the supplied value placed in the low-order indices.
func ToBytes48(val []byte) [48]byte {
	var res [48]byte
	copy(res[:], val)
	return res
}

// ToBytes64 returns a 64-byte array with the supplied value placed in the low-order indices.
func ToBytes64(val []byte) [64]byte {
	var res [64]byte
	copy(res[:], val)
	return res
}

// ToBytes96 returns a 96-byte array with the supplied value placed in the low-order indices.
func ToBytes96(val []byte) [96]byte {
	var res [96]byte
	copy(res[:], val)
	return res
}

// FromHexString returns a byte array given a hex string
func FromHexString(data string) ([]byte, error) {
	data = strings.TrimPrefix(data, "0x")
	if len(data)%2 == 1 {
		// Odd number of characters; even it up
		data = "0" + data
	}
	return hex.DecodeString(data)
}

// XOR returns an XORd copy of the bytes.
func XOR(data []byte) []byte {
	res := make([]byte, len(data))
	for i := range data {
		res[i] = data[i] ^ 0xff
	}
	return res
}
