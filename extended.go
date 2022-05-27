// Package extended provides conversions to and from 80-bit "extended"
// floating-point numbers.
//
// Note that while NaNs are handled by this package, the distinction between
// quiet NaN and signaling NaN is not preserved.
package extended

import (
	"math"
	"math/bits"
)

// An Extended is an 80-bit extended precision floating-point number.
type Extended struct {
	// The sign is stored as the high bit, the low 15 bits contain the exponent,
	// with a bias of 16383.
	SignExponent uint16

	// The fraction includes a ones place as the high bit. The valiue in the
	// ones place may be zero.
	Fraction uint64
}

// FromFloat64 converts a 64-bit floating-point number to an 80-bit extended
// floating-point number.
func FromFloat64(x float64) (e Extended) {
	xbits := math.Float64bits(x)
	sign := int(xbits>>(63-15)) & 0x8000
	exponent := int(xbits>>52) & ((1 << 11) - 1)
	mantissa := xbits & ((1 << 52) - 1)
	switch exponent {
	case 0:
		// Zero or subnormal.
		// Number is (-1)^sign * 2^-1022 * 0.mantissa.
		if mantissa == 0 {
			return Extended{uint16(sign), 0}
		}
		// 2^-1022 * 0.mantissa = 2^(e-16383) * 2^lzero * 0.mantissa
		// -1022 = e - 16383 + lzero
		// e = -1022 + 16383 - lzero
		nzero := bits.LeadingZeros64(mantissa)
		exponent := 16383 - 1022 + 11 - nzero
		return Extended{uint16(sign | exponent), mantissa << nzero}

	case (1 << 11) - 1:
		// Infinity or NaN.
		if mantissa == 0 {
			return Extended{uint16(sign | 0x7fff), 0}
		}
		return Extended{uint16(sign | 0x7fff), ^uint64(0)}

	default:
		// 2^(e64 - 1023) * 1.fraction = 2^(e80 - 16383) * 1.fraction
		// e63 - 1023 = e80 - 16383
		// e80 = e63 + 16383 - 1023
		exponent := exponent + 16383 - 1023
		return Extended{uint16(sign | exponent), 1<<63 | mantissa<<11}
	}
}
