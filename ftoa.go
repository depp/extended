package extended

import "fmt"

// String converts the extended-precision value to a string.
func (e Extended) String() string {
	return e.Text('g', 10)
}

// Text converts the floating-point number to a string using the given format
// specifier and precision.
func (e Extended) Text(format byte, prec int) string {
	cap := 10
	if prec > 0 {
		cap += prec
	}
	return string(e.Append(make([]byte, 0, cap), format, prec))
}

// Append appends the string form of the number to buf and returns the result.
func (e Extended) Append(buf []byte, fmt byte, prec int) []byte {
	if f, err := e.BigFloat(); err == nil {
		return f.Append(buf, fmt, prec)
	}
	return append(buf, "NaN"...)
}

var space = []byte{' '}

func writeSpace(s fmt.State, n int) {
	for i := 0; i < n; i++ {
		s.Write(space)
	}
}

// Format implements fmt.Formatter.
func (e Extended) Format(s fmt.State, format rune) {
	if f, err := e.BigFloat(); err == nil {
		f.Format(s, format)
		return
	}
	var padding int
	if width, hasWidth := s.Width(); hasWidth && width > 3 {
		padding = width - 3
	}
	if s.Flag('-') {
		s.Write([]byte("NaN"))
		writeSpace(s, padding)
	} else {
		writeSpace(s, padding)
		s.Write([]byte("NaN"))
	}
}
