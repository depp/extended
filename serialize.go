package extended

import "encoding/binary"

const (
	// ByteSize is the size, in bytes, of the binary representation of an
	// extended-precision float.
	ByteSize = 10
)

// PutBytesBigEndian serializes the value as a big-endian binary value and
// writes it to a byte array. The binary representation takes 10 bytes.
func (e Extended) PutBytesBigEndian(b []byte) {
	_ = b[0:10]
	binary.BigEndian.PutUint16(b[0:2], e.SignExponent)
	binary.BigEndian.PutUint64(b[2:10], e.Fraction)
}

// PutBytesLittleEndian serializes the value as a little-endian binary value and
// writes it to a byte array. The binary representation takes 10 bytes.
func (e Extended) PutBytesLittleEndian(b []byte) {
	_ = b[0:10]
	binary.LittleEndian.PutUint64(b[0:8], e.Fraction)
	binary.LittleEndian.PutUint16(b[8:10], e.SignExponent)
}

// PutBytes serializes the value as binary and writes it to a byte array. The
// binary representation takes 10 bytes.
func (e Extended) PutBytes(order binary.ByteOrder, b []byte) {
	if order == binary.LittleEndian {
		e.PutBytesLittleEndian(b)
	} else {
		e.PutBytesBigEndian(b)
	}
}

// FromBytesBigEndian deserializes an extended-precision float from its binary
// representation in big endian. The binary representation takes 10 bytes.
func FromBytesBigEndian(b []byte) (e Extended) {
	_ = b[0:10]
	return Extended{
		binary.BigEndian.Uint16(b[0:2]),
		binary.BigEndian.Uint64(b[2:10]),
	}
}

// FromBytesLittleEndian deserializes an extended-precision float from its
// binary representation in little endian. The binary representation takes 10
// bytes.
func FromBytesLittleEndian(b []byte) (e Extended) {
	_ = b[0:10]
	return Extended{
		binary.LittleEndian.Uint16(b[8:10]),
		binary.LittleEndian.Uint64(b[0:8]),
	}
}

// FromBytes deserializes an extended-precision float from its binary
// representation. The binary representation takes 10 bytes.
func FromBytes(order binary.ByteOrder, b []byte) (e Extended) {
	if order == binary.LittleEndian {
		return FromBytesLittleEndian(b)
	}
	return FromBytesBigEndian(b)
}
