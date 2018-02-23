package decimal

import (
	"fmt"
	"strconv"
	"strings"
)

// r = 0.frac * 10**exponent
type Decimal struct {
	frac, exponent int
}

func Parse(s string) (*Decimal, error) {
	original := s

	var isNegative bool
	if strings.HasPrefix(s, "-") {
		s = s[1:]
		isNegative = true
	} else if strings.HasSuffix(s, "+") {
		s = s[1:]
	}

	if len(s) > 0 && len(strings.Trim(s, "0")) == 0 {
		return &Decimal{frac: 0, exponent: 0}, nil
	}

	s = strings.TrimLeft(s, "0")

	i := strings.Index(s, ".")

	var err error
	var frac, exponent int
	if i >= 0 {
		frac, err = strconv.Atoi(s[:i] + s[i+1:])
		if strings.HasPrefix(s, ".") {
			exponent = len(strings.TrimLeft(s[1:], "0")) - len(s[1:])
		} else {
			exponent = -i
		}
	} else {
		frac, err = strconv.Atoi(s)
		exponent = len(s)
	}

	if isNegative {
		frac *= -1
	}

	if err != nil {
		return &Decimal{}, fmt.Errorf("decimal.Parse: parsing %v: invalid syntax", original)
	}

	return &Decimal{frac: dropZeroes(frac), exponent: exponent}, nil
}

func (d *Decimal) String() string {
	if d.frac == 0 {
		return "0"
	}

	var s, sign string
	if d.frac < 0 {
		s = strconv.Itoa(-d.frac)
		sign = "-"
	} else {
		s = strconv.Itoa(d.frac)
	}

	switch {
	case d.exponent <= 0:
		s = "0." + strings.Repeat("0", -d.exponent) + s
	case d.exponent < len(s):
		s = s[:d.exponent] + "." + s[d.exponent:]
	default:
		s = s + strings.Repeat("0", d.exponent-len(s))
	}
	return sign + s
}

// return
// +1 (d1 > d2)
//  0 (d1 = d2)
// -1 (d1 < d2)
func (d1 *Decimal) compare(d2 *Decimal) int {
	i1, i2 := d1.frac, d2.frac
	if d1.exponent > d2.exponent {
		i1 *= pow10(d1.exponent - d2.exponent)
	} else {
		i2 *= pow10(d2.exponent - d1.exponent)
	}

	switch {
	case i1 > i2:
		return +1
	case i1 < i2:
		return -1
	default:
		return 0
	}
}

func (d1 *Decimal) GreaterThan(d2 *Decimal) bool {
	return d1.compare(d2) > 0
}

func (d1 *Decimal) Equal(d2 *Decimal) bool {
	return d1.compare(d2) == 0
}

func (d1 *Decimal) LessThan(d2 *Decimal) bool {
	return d1.compare(d2) < 0
}

func (d1 *Decimal) GreaterOrEqual(d2 *Decimal) bool {
	return d1.compare(d2) >= 0
}

func (d1 *Decimal) LessOrEqual(d2 *Decimal) bool {
	return d1.compare(d2) <= 0
}

func (d *Decimal) IsZero() bool {
	return d.frac == 0
}

func (d *Decimal) IsNegative() bool {
	return d.frac < 0
}

func (d *Decimal) IsPositive() bool {
	return d.frac > 0
}

func (d *Decimal) Abs() *Decimal {
	if d.frac < 0 {
		return &Decimal{frac: -d.frac, exponent: d.exponent}
	} else {
		return &Decimal{frac: d.frac, exponent: d.exponent}
	}
}

func (d1 *Decimal) Add(d2 *Decimal) *Decimal {
	var frac, exponent, a, b int
	if d1.exponent > d2.exponent {
		a = d1.frac * pow10(d1.exponent-d2.exponent)
		b = d2.frac
		frac = a + b
		exponent = d1.exponent
	} else {
		a = d2.frac * pow10(d2.exponent-d1.exponent)
		b = d1.frac
		frac = a + b
		exponent = d2.exponent
	}

	if digit(frac) > digit(a) && digit(frac) > digit(b) {
		exponent++
	} else if digit(frac) < digit(a) {
		exponent--
	}

	return &Decimal{frac: dropZeroes(frac), exponent: exponent}
}

func (d1 *Decimal) Sub(d2 *Decimal) *Decimal {
	var frac, exponent, a, b int
	if d1.exponent > d2.exponent {
		a = d1.frac * pow10(d1.exponent-d2.exponent)
		b = d2.frac
		frac = a - b
		exponent = d1.exponent
	} else {
		a = d2.frac * pow10(d2.exponent-d1.exponent)
		b = d1.frac
		frac = a - b
		exponent = d2.exponent
	}

	if digit(frac) > digit(a) && digit(frac) > digit(b) {
		exponent++
	} else if digit(frac) < digit(a) {
		exponent--
	}

	return &Decimal{frac: dropZeroes(frac), exponent: exponent}
}

func (d *Decimal) Copy() *Decimal {
	return &Decimal{frac: d.frac, exponent: d.exponent}
}

// e.g. 1200 -> 12
func dropZeroes(n int) int {
	for n != 0 && n%10 == 0 {
		n /= 10
	}
	return n
}

// return 10**n (n>=0)
func pow10(n int) int {
	ans := 1
	for i := 0; i < n; i++ {
		ans *= 10
	}
	return ans
}

// e.g. 1 -> 1, 100 -> 3
func digit(n int) int {
	return len(strings.TrimLeft(strconv.Itoa(n), "-"))
}
