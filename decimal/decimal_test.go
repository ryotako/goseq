package decimal

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{input: "1", want: "1"},
		{input: "10", want: "10"},
		{input: "10000", want: "10000"},
		{input: "0.1", want: "0.1"},
		{input: "0.0001", want: "0.0001"},
	}

	for _, test := range tests {
		got, err := Parse(test.input)
		if err != nil {
			t.Errorf("%q >> parse error", test.input)
		} else if got.String() != test.want {
			t.Errorf("%q >> want %q, but %q(%#v)", test.input, test.want, got.String(), got)
		}
	}
}

func TestCondition(t *testing.T) {
	tests := []struct {
		msg  string
		cond bool
	}{
		{msg: "1 = 1", cond: D("1").Equal(D("1"))},
		{msg: "0.001 = 0.001", cond: D("0.001").Equal(D("0.001"))},
		{msg: "1000 = 1000", cond: D("1000").Equal(D("1000"))},
		{msg: "1 < 2", cond: D("1").LessThan(D("2"))},
		{msg: "0.001 < 0.1", cond: D("0.001").LessThan(D("0.1"))},
		{msg: "10 < 1000", cond: D("10").LessThan(D("1000"))},
		{msg: "2 > 1", cond: D("2").GreaterThan(D("1"))},
		{msg: "0.1 > 0.001", cond: D("0.1").GreaterThan(D("0.001"))},
		{msg: "1000 > 10", cond: D("1000").GreaterThan(D("10"))},
		{msg: "0 is zero", cond: D("0").IsZero()},
		{msg: "+000 is zero", cond: D("000").IsZero()},
		{msg: "-0.0 is zero", cond: D("000").IsZero()},
		{msg: "0.1 is positive", cond: D("0.1").IsPositive()},
		{msg: "-0.1 is negative", cond: D("-0.1").IsNegative()},
	}

	for _, test := range tests {
		if !test.cond {
			t.Errorf("%q >> false", test.msg)
		}
	}
}

func TestCalc(t *testing.T) {
	tests := []struct {
		msg, want string
		got       *Decimal
	}{
		{msg: "|1|", got: D("1").Abs(), want: "1"},
		{msg: "|-1|", got: D("-1").Abs(), want: "1"},
		{msg: "|0|", got: D("0").Abs(), want: "0"},
		{msg: "|0.001|", got: D("0.001").Abs(), want: "0.001"},
		{msg: "|-0.001|", got: D("-0.001").Abs(), want: "0.001"},
		{msg: "1+1", got: D("1").Add(D("1")), want: "2"},
		{msg: "3+7", got: D("3").Add(D("7")), want: "10"},
		{msg: "10+2", got: D("10").Add(D("2")), want: "12"},
		{msg: "1+0.2", got: D("1").Add(D("0.2")), want: "1.2"},
		{msg: "1-1", got: D("1").Sub(D("1")), want: "0"},
		{msg: "10-1", got: D("10").Sub(D("1")), want: "9"},
		{msg: "1-0.1", got: D("1").Sub(D("0.1")), want: "0.9"},
	}

	for _, test := range tests {
		if test.got.String() != test.want {
			t.Errorf("%q >> want %q, but %q(%#v)", test.msg, test.want, test.got.String(), test.got)
		}
	}
}

func D(s string) *Decimal {
	d, err := Parse(s)
	if err != nil {
		panic(s + ": parse error")
	}
	return d
}
