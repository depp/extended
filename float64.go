package extended

import (
	"math"
	"math/bits"
)

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

// Float64 returns the value of this 80-bit floating-point number as a float64.
// The result is rounded to the nearest float64, breaking ties towards even in
// the least-significant bit. Values which, after rounding, would be outside the
// range of a float64 are flushed to zero or infinity.
func (e Extended) Float64() float64 {
	exponent := int(e.SignExponent) & 0x7fff
	sign := int(e.SignExponent) & 0x8000
	if exponent == 0x7fff {
		if e.Fraction == 0 {
			return math.Inf(-sign)
		}
		return math.Copysign(math.NaN(), float64(-sign))
	}
	if e.Fraction == 0 {
		return math.Copysign(0, float64(-sign))
	}
	// 2^(e64 - 1023) * 1.fraction
	// = 2^(e80 - 16383) * 1.fraction / 2^nzero
	// e63 - 1023 = e80 - 16383
	// e63 = e80 - 16383 + 1023 - nzero
	nzero := bits.LeadingZeros64(e.Fraction)
	exponent += 1023 - 16383 - nzero
	fraction := e.Fraction << nzero
	if exponent <= 0 {
		// Subnormal numbers.
		shift := 12 - exponent
		var rem uint64
		if shift < 64 {
			fraction, rem = fraction>>shift, fraction<<(64-shift)
		} else if shift == 64 {
			fraction, rem = 0, fraction
		} else {
			fraction = 0
		}
		// The (fraction & 1) makes this round to even.
		if rem|(fraction&1) > 1<<63 {
			fraction++
		}
		exponent = 0
	} else {
		// Round to 52 bits. The addition of ((fraction >> 11) & 1) makes this
		// round to even.
		rem := fraction&((1<<11)-1) | (fraction>>11)&1
		fraction = (fraction >> 11) & ((1 << 52) - 1)
		if rem > 1<<10 {
			if fraction < (1<<52)-1 {
				fraction++
			} else {
				exponent++
				fraction = 0
			}
		}
		if exponent >= (1<<11)-1 {
			return math.Inf(-sign)
		}
	}
	return math.Float64frombits(
		(uint64(e.SignExponent)&0x8000)<<48 |
			uint64(exponent)<<52 |
			fraction)
}
