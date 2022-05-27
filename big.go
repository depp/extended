package extended

import (
	"errors"
	"math/big"
)

// ErrIsNaN indicates that the value is NaN and cannot be converted.
var ErrIsNaN = errors.New("value is NaN")

// BigFloat converts the number to an arbitrary-precision float.
func (e Extended) BigFloat() (*big.Float, error) {
	signbit := e.SignExponent&0x8000 != 0
	exponent := int(e.SignExponent) & 0x7fff
	if exponent == 0x7fff {
		if e.Fraction == 0 {
			return new(big.Float).SetInf(signbit), nil
		}
		return nil, ErrIsNaN
	}
	f := new(big.Float)
	f.SetUint64(e.Fraction)
	if signbit {
		f.Neg(f)
	}
	f.SetMantExp(f, exponent-16383-63)
	return f, nil
}
