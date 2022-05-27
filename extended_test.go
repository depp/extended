package extended

import (
	"math"
	"strconv"
	"testing"
)

func TestFromF64(t *testing.T) {
	type testcase struct {
		name     string
		exponent uint16
		fraction uint64
		input    float64
	}
	cases := []testcase{
		{"basic_one", 16383, 1 << 63, 1.0},
		{"basic_two", 16384, 1 << 63, 2.0},
		{"basic_half", 16382, 1 << 63, 0.5},
		{"small", 16383 - 10, 1 << 63, 0.0009765625},
		{"smaller", 16383 - 100, 1 << 63, 7.888609052210118e-31},
		{"after_one", 16383, (1 << 63) + (1 << 11), 1.0000000000000002},
		{"infinity", 32767, 0, math.Inf(0)},
		{"zero", 0, 0, 0.0},
		{"nan", 32767, ^uint64(0), math.NaN()},
		{"smallest_normal", 15361, 1 << 63, 2.2250738585072014e-308},
		{"subnormal", 15360, 1 << 63, 1.1125369292536007e-308},
		{"smallest_subnormal", 15309, 1 << 63, 5e-324},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for sign := 0; sign < 2; sign++ {
				expect := Extended{
					c.exponent | uint16(sign<<15),
					c.fraction,
				}
				fin := c.input
				if sign != 0 {
					fin = -fin
				}
				out := FromFloat64(fin)
				if out != expect {
					t.Errorf("FromFloat64(%s) = %04x:%016x, expect %04x:%016x",
						strconv.FormatFloat(fin, 'g', -1, 64),
						out.SignExponent, out.Fraction, expect.SignExponent, expect.Fraction)
				}
			}
		})
	}
}
