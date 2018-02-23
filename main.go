package main

// TODO: 浮動小数の誤差の処理
// 現状では goseq 1 0.1 2 が 1, 1.1 ... 1.9 で止まる

import (
	"fmt"
	"github.com/shopspring/decimal"
	"io"
	"os"
)

type CLI struct {
	outStream, errStream io.Writer
}

func main() {
	cli := &CLI{outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(cli.Run(os.Args[1:]))
}

func usage(w io.Writer) {
	fmt.Fprintln(w, "usage: seq [first [incr]] last")
}

func (c *CLI) Run(args []string) int {

	var fst, inc, lst decimal.Decimal
	errs := make([]error, 3)

	switch len(args) {
	case 1:
		fst, _ = decimal.NewFromString("1")
		inc, _ = decimal.NewFromString("1")
		lst, errs[0] = decimal.NewFromString(args[0])
	case 2:
		fst, errs[0] = decimal.NewFromString(args[0])
		inc, _ = decimal.NewFromString("1")
		lst, errs[1] = decimal.NewFromString(args[1])
	case 3:
		fst, errs[0] = decimal.NewFromString(args[0])
		inc, errs[1] = decimal.NewFromString(args[1])
		lst, errs[2] = decimal.NewFromString(args[2])
	default:
		usage(c.errStream)
		return 1
	}

	for i, err := range errs {
		if err != nil {
			errorf(c.errStream, "invalid floating point argument: %s\n", args[i])
			return 1
		}
	}

	if len(args) < 3 && fst.GreaterThan(lst) {
		inc, _ = decimal.NewFromString("-1")
	}

	if fst.LessThan(lst) {
		switch inc.Sign() {
		case 0:
			errorf(c.errStream, "zero increment")
			return 1
		case -1:
			errorf(c.errStream, "needs positive increment")
			return 1
		default:
			for i := fst; i.LessThanOrEqual(lst); i = i.Add(inc) {
				fmt.Fprintln(c.outStream, i)
			}
		}
	} else {
		switch inc.Sign() {
		case 0:
			errorf(c.errStream, "zero increment")
			return 1
		case 1:
			errorf(c.errStream, "needs negative increment")
			return 1
		default:
			for i := fst; i.GreaterThanOrEqual(lst); i = i.Add(inc) {
				fmt.Fprintln(c.outStream, i)
			}
		}
	}
	return 0
}

func errorf(w io.Writer, format string, a ...interface{}) {
	fmt.Fprintf(w, fmt.Sprintf("goseq: %s", format), a...)
}
