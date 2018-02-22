package main

// TODO: 浮動小数の誤差の処理
// 現状では goseq 1 0.1 2 が 1, 1.1 ... 1.9 で止まる

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"io"
	"os"
	"strconv"
	"strings"
)

type CLI struct {
	outStream, errStream io.Writer
}

// structure to calc floating point number correctly
// 1.23 => Number{number:123, order:2}
// 0.001 => Number{number:1, order:3}
type Number struct {
	number, order int
}

func ParseNumber(s string) (*Number, error) {
	if len(strings.Trim(s, "0")) == 0 {
		return &Number{number: 0, order: 0}, nil
	}

	var isNegative bool
	if strings.HasPrefix(s, "+") {
		s = s[1:]
	} else if strings.HasPrefix(s, "-") {
		isNegative = true
		s = s[1:]
	}
	s = strings.Trim(s, "0")
	i := strings.Index(s, ".")

	var order int
	if i < 0 {
		order = 0
	} else {
		order = len(s) - i - 1
	}

	number, err := strconv.Atoi(strings.Replace(s, ".", "", 1))
	if isNegative {
		number *= -1
	}

	if err != nil {
		return &Number{}, fmt.Errorf("Failed to parse %v as a number", s)
	} else {
		return &Number{number: number, order: order}, nil
	}
}

func (n *Number) String() string {
	number := n.number
	sign := ""
	if number < 0 {
		number *= -1
		sign = "-"
	}

	s := strconv.Itoa(number)
	l := len(s)
	switch {
	case n.order < 0:
		return ""
	case n.order == 0:
		return sign + s
	case n.order >= l:
		return sign + "0." + strings.Repeat("0", n.order-l) + s
	default:
		return fmt.Sprintf("%s%s.%s", sign, s[:l-n.order], s[l-n.order:])
	}
}

func (n *Number) Copy() *Number {
	return &Number{number: n.number, order: n.order}
}

func (n1 *Number) Add(n2 *Number) *Number {
	var order, number int
	if n1.order > n2.order {
		order = n1.order
		number = n1.number + n2.number*iPow(10, n1.order-n2.order)
	} else {
		order = n2.order
		number = n2.number + n1.number*iPow(10, n2.order-n1.order)
	}

	for number > 0 && number%10 == 0 {
		number /= 10
		order -= 1
	}

	return &Number{number: number, order: order}
}

func (n1 *Number) Compare(n2 *Number) int {
	var i1, i2 int
	if n1.order > n2.order {
		i1 = n1.number
		i2 = n2.number * iPow(10, n1.order-n2.order)
	} else {
		i1 = n1.number * iPow(10, n2.order-n1.order)
		i2 = n2.number
	}

	switch {
	case i1 > i2:
		return 1
	case i1 < i2:
		return -1
	default:
		return 0
	}
}

func (n1 *Number) GreaterThan(n2 *Number) bool {
	return n1.Compare(n2) > 0
}

func (n1 *Number) Equal(n2 *Number) bool {
	return n1.Compare(n2) == 0
}

func (n1 *Number) LessThan(n2 *Number) bool {
	return n1.Compare(n2) < 0
}

func (n1 *Number) GreaterOrEqual(n2 *Number) bool {
	return n1.Compare(n2) >= 0
}

func (n1 *Number) LessOrEqual(n2 *Number) bool {
	return n1.Compare(n2) <= 0
}

func (n *Number) IsZero() bool {
	return n.number == 0
}

func (n *Number) IsNegative() bool {
	return n.number < 0
}

func (n *Number) IsPositive() bool {
	return n.number > 0
}

func main() {
	cli := &CLI{outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(cli.Run(os.Args[1:]))
}

func usage(w io.Writer) {
	fmt.Fprintln(w, "usage: seq [first [incr]] last")
}

func printNum(w io.Writer, n float64) {
	fmt.Fprintln(w, humanize.Ftoa(n))
}

func (c *CLI) Run(args []string) int {

	var fst, inc, lst *Number
	var err0, err1, err2 error

	switch len(args) {
	case 1:
		fst, _ = ParseNumber("1")
		inc, _ = ParseNumber("1")
		lst, err0 = ParseNumber(args[0])
	case 2:
		fst, err0 = ParseNumber(args[0])
		inc, _ = ParseNumber("1")
		lst, err1 = ParseNumber(args[1])
	case 3:
		fst, err0 = ParseNumber(args[0])
		inc, err1 = ParseNumber(args[1])
		lst, err2 = ParseNumber(args[2])
	default:
		usage(c.errStream)
		return 1
	}

	if err0 != nil {
		fmt.Fprintf(c.errStream, "invalid floating point argument: %s\n", args[0])
		return 1
	}
	if err1 != nil {
		fmt.Fprintf(c.errStream, "invalid floating point argument: %s\n", args[1])
		return 1
	}
	if err2 != nil {
		fmt.Fprintf(c.errStream, "invalid floating point argument: %s\n", args[2])
		return 1
	}

	if len(args) < 3 && fst.GreaterThan(lst) {
		inc, _ = ParseNumber("-1")
	}

	if fst.LessThan(lst) {
		if inc.IsZero() {
			fmt.Fprintln(c.errStream, "zero increment")
			return 1
		} else if inc.IsNegative() {
			fmt.Fprintln(c.errStream, "needs positive increment")
			return 1
		} else {
			for i := fst.Copy(); i.LessOrEqual(lst); i = i.Add(inc) {
				fmt.Fprintln(c.outStream, i)
			}
		}
	} else {
		if inc.IsZero() {
			fmt.Fprintln(c.errStream, "zero increment")
			return 1
		} else if fst.Equal(lst) {
			fmt.Fprintln(c.outStream, fst)
		} else if inc.IsPositive() {
			fmt.Fprintln(c.errStream, "needs negative increment")
			return 1
		} else {
			for i := fst.Copy(); i.GreaterOrEqual(lst); i = i.Add(inc) {
				fmt.Fprintln(c.outStream, i)
			}
		}
	}
	return 0
}

// return x**y
func iPow(x, y int) int {
	if y < 0 {
		if x == 0 {
			return 1
		} else {
			return 0
		}
	}

	ans := 1
	for i := 0; i < y; i++ {
		ans *= x
	}
	return ans
}
