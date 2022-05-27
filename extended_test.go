package extended

import (
	"bytes"
	"encoding/binary"
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

func equal(x, y float64) bool {
	return math.Float64bits(x) == math.Float64bits(y)
}

func TestToFloat64(t *testing.T) {
	type testcase struct {
		name     string
		exponent uint16
		fraction uint64
		output   float64
	}
	cases := []testcase{
		{"basic_one", 16383, 1 << 63, 1.0},
		{"basic_two", 16384, 1 << 63, 2.0},
		{"basic_half", 16382, 1 << 63, 0.5},
		{"after_one", 16383, (1 << 63) + (1 << 11), 1.0000000000000002},
		{"round_even_1", 16383, (1 << 63) + (1 << 10), 1.0},
		{"round_even_2", 16383, (1 << 63) + (1 << 10) + 1, 1.0000000000000002},
		{"round_even_3", 16383, (1 << 63) + (1 << 11), 1.0000000000000002},
		{"round_even_4", 16383, (1 << 63) + (3 << 10) - 1, 1.0000000000000002},
		{"round_even_5", 16383, (1 << 63) + (3 << 10), 1.0000000000000004},
		{"round_exponent", 16381, ^uint64(0), 0.5},
		{"inf", 32767, 0, math.Inf(0)},
		{"large_1", 32000, 1 << 63, math.Inf(0)},
		{"large_2", 32000, ^uint64(0), math.Inf(0)},
		{"large_3", 17406, 0xfffffffffffff800, 1.7976931348623157e+308},
		{"large_4", 17406, 0xfffffffffffffbff, 1.7976931348623157e+308},
		{"large_5", 17406, 0xfffffffffffffc00, math.Inf(0)},
		{"zero", 0, 0, 0.0},
		{"nan_1", 32767, 1, math.NaN()},
		{"nan_2", 32767, 1 << 63, math.NaN()},
		{"smallest_normal", 15361, 1 << 63, 2.2250738585072014e-308},
		{"subnormal", 15360, 1 << 63, 1.1125369292536007e-308},
		{"smallest_subnormal", 15309, 1 << 63, 5e-324},
		{"smallest_subnormal_roundup", 15308, (1 << 63) + 1, 5e-324},
		{"smallest_subnormal_rounddown", 15308, 1 << 63, 0.0},
		{"round_to_zero", 10000, 1 << 63, 0.0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for sign := 0; sign < 2; sign++ {
				in := Extended{c.exponent, c.fraction}
				expect := c.output
				if sign == 1 {
					in.SignExponent |= 0x8000
					expect = -expect
				}
				out := in.Float64()
				if !equal(out, expect) {
					t.Errorf("%04x:%016x.Float64() = %s, expect %s",
						in.SignExponent, in.Fraction,
						strconv.FormatFloat(out, 'g', -1, 64),
						strconv.FormatFloat(expect, 'g', -1, 64))
					t.Logf("bits: %016x %016x", math.Float64bits(out), math.Float64bits(expect))
				}
				if big, err := in.BigFloat(); err == nil {
					out2, _ := big.Float64()
					if !equal(out2, out) {
						t.Errorf("%04x:%016x.BigFloat().Float64() = %s, expect %s",
							in.SignExponent, in.Fraction,
							strconv.FormatFloat(out2, 'g', -1, 64),
							strconv.FormatFloat(out, 'g', -1, 64))
						t.Logf("bits: %016x %016x", math.Float64bits(out), math.Float64bits(expect))
					}
				}
			}
		})
	}
}

func TestSerialize(t *testing.T) {
	var (
		e  = Extended{0x1234, 0x1122334455667788}
		be = []byte{0x12, 0x34, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88}
		le = []byte{0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0x34, 0x12}
	)

	if out := FromBytesBigEndian(be); out != e {
		t.Errorf("FromBytesBigEndian = %v, expect %v", out, e)
	}
	if out := FromBytes(binary.BigEndian, be); out != e {
		t.Errorf("FromBytes(BigEndian) = %v, expect %v", out, e)
	}
	if out := FromBytesLittleEndian(le); out != e {
		t.Errorf("FromBytesLittleEndian = %v, expect %v", out, e)
	}
	if out := FromBytes(binary.LittleEndian, le); out != e {
		t.Errorf("FromBytes(LittleEndian) = %v, expect %v", out, e)
	}

	{
		var b [10]byte
		e.PutBytesBigEndian(b[:])
		if !bytes.Equal(b[:], be) {
			t.Errorf("PutBytesBigEndian = %v, expect %v", b[:], be)
		}
	}
	{
		var b [10]byte
		e.PutBytesLittleEndian(b[:])
		if !bytes.Equal(b[:], le) {
			t.Errorf("PutBytesLittleEndian = %v, expect %v", b[:], le)
		}
	}
	{
		var b [10]byte
		e.PutBytes(binary.BigEndian, b[:])
		if !bytes.Equal(b[:], be) {
			t.Errorf("PutBytes(BigEndian) = %v, expect %v", b[:], be)
		}
	}
	{
		var b [10]byte
		e.PutBytes(binary.LittleEndian, b[:])
		if !bytes.Equal(b[:], le) {
			t.Errorf("PutBytes(LittleEndian) = %v, expect %v", b[:], le)
		}
	}
}
