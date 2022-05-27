// Package extended provides conversions to and from 80-bit "extended"
// floating-point numbers (float80).
//
// Note that while NaNs are handled by this package, the distinction between
// quiet NaN and signaling NaN is not preserved during conversions.
package extended

// An Extended is an 80-bit extended precision floating-point number.
type Extended struct {
	// The sign is stored as the high bit, the low 15 bits contain the exponent,
	// with a bias of 16383.
	SignExponent uint16

	// The fraction includes a ones place as the high bit. The valiue in the
	// ones place may be zero.
	Fraction uint64
}
