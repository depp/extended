package main

import (
	"encoding/binary"
	"fmt"

	"github.com/depp/extended"
)

func main() {
	e := extended.Extended{
		SignExponent: 0x3fff,
		Fraction:     0xC000000000000000,
	}
	// Value: 1.500
	fmt.Printf("Value: %.3f\n", e)

	// Float64: 1.500
	f64 := e.Float64()
	fmt.Printf("Float64: %.3f\n", f64)

	// Value: 100.75
	// SignExponent: 0x4005
	// Fraction: 0xc980000000000000
	e = extended.FromFloat64(100.75)
	fmt.Println("Value:", e)
	fmt.Printf("SignExponent: 0x%04x\n", e.SignExponent)
	fmt.Printf("Fraction: 0x%016x\n", e.Fraction)

	// Binary (big endian): 4005c980000000000000
	var buf [extended.ByteSize]byte
	e.PutBytes(binary.BigEndian, buf[:])
	fmt.Printf("Binary (big endian): %x\n", buf[:])

	// Binary (little endian): 00000000000080c90540
	e.PutBytes(binary.LittleEndian, buf[:])
	fmt.Printf("Binary (little endian): %x\n", buf[:])
}
