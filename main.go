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
	var err0, err1, err2 error

	switch len(args) {
	case 1:
		fst, _ = decimal.NewFromString("1")
		inc, _ = decimal.NewFromString("1")
		lst, err0 = decimal.NewFromString(args[0])
	case 2:
		fst, err0 = decimal.NewFromString(args[0])
		inc, _ = decimal.NewFromString("1")
		lst, err1 = decimal.NewFromString(args[1])
	case 3:
		fst, err0 = decimal.NewFromString(args[0])
		inc, err1 = decimal.NewFromString(args[1])
		lst, err2 = decimal.NewFromString(args[2])
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
		inc, _ = decimal.NewFromString("-1")
	}

	if fst.LessThan(lst) {
		switch inc.Sign() {
		case 0:
			fmt.Fprintln(c.errStream, "zero increment")
			return 1
		case -1:
			fmt.Fprintln(c.errStream, "needs positive increment")
			return 1
		default:
			for i := fst; i.LessThanOrEqual(lst); i = i.Add(inc) {
				fmt.Fprintln(c.outStream, i)
			}
		}
	} else {
		switch inc.Sign() {
		case 0:
			fmt.Fprintln(c.errStream, "zero increment")
			return 1
		case 1:
			fmt.Fprintln(c.errStream, "needs negative increment")
			return 1
		default:
			for i := fst; i.GreaterThanOrEqual(lst); i = i.Add(inc) {
				fmt.Fprintln(c.outStream, i)
			}
		}
	}
	return 0
}
