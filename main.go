package main

import (
	"fmt"
	"github.com/shopspring/decimal"
	"io"
	"os"
	"regexp"
	"strings"
)

const (
	SUCCESS = iota
	INVALID_SYNTAX
	INCREMENT_ERROR
)

type CLI struct {
	outStream, errStream io.Writer
}

func main() {
	cli := &CLI{outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(cli.Run(os.Args[1:]))
}

func usage(w io.Writer) {
	fmt.Fprintln(w, "usage: goseq [first [incr]] last")
}

func (c *CLI) Run(args []string) int {
	flagS := "\n"
	flagT := ""
	numArgs := []string{}

loop:
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-s":
			if i+1 < len(args) {
				flagS = args[i+1]
				i++
			} else {
				errorf(c.errStream, "option requires an argument -- s")
				return INVALID_SYNTAX
			}

		case "-t":
			if i+1 < len(args) {
				flagT = args[i+1]
				i++
			} else {
				errorf(c.errStream, "option requires an argument -- s")
				return INVALID_SYNTAX
			}

		default:
			if r := regexp.MustCompile(`^-[^\d\.]`); r.MatchString(args[i]) {
				errorf(c.errStream, "illegal option -- %s", strings.TrimLeft(args[i], "-"))
				return INVALID_SYNTAX
			} else {
				numArgs = append(numArgs, args[i:]...)
				break loop
			}
		}
	}

	var fst, inc, lst decimal.Decimal
	errs := make([]error, 3)

	switch len(numArgs) {
	case 1:
		fst, _ = decimal.NewFromString("1")
		inc, _ = decimal.NewFromString("1")
		lst, errs[0] = decimal.NewFromString(numArgs[0])
	case 2:
		fst, errs[0] = decimal.NewFromString(numArgs[0])
		inc, _ = decimal.NewFromString("1")
		lst, errs[1] = decimal.NewFromString(numArgs[1])
	case 3:
		fst, errs[0] = decimal.NewFromString(numArgs[0])
		inc, errs[1] = decimal.NewFromString(numArgs[1])
		lst, errs[2] = decimal.NewFromString(numArgs[2])
	default:
		usage(c.errStream)
		return INVALID_SYNTAX
	}

	for i, err := range errs {
		if err != nil {
			errorf(c.errStream, "invalid floating point argument: %s\n", numArgs[i])
			return INVALID_SYNTAX
		}
	}

	if len(numArgs) < 3 && fst.GreaterThan(lst) {
		inc, _ = decimal.NewFromString("-1")
	}

	if fst.LessThan(lst) {
		switch inc.Sign() {
		case 0:
			errorf(c.errStream, "zero increment")
			return INCREMENT_ERROR
		case -1:
			errorf(c.errStream, "needs positive increment")
			return INCREMENT_ERROR
		default:
			for i := fst; i.LessThanOrEqual(lst); i = i.Add(inc) {
				f, _ := i.Float64()
				fmt.Fprintf(c.outStream, "%s%s", fmt.Sprintf(flagF, f), flagS)
			}
		}
	} else if fst.GreaterThan(lst) {
		switch inc.Sign() {
		case 0:
			errorf(c.errStream, "zero increment")
			return INCREMENT_ERROR
		case 1:
			errorf(c.errStream, "needs negative increment")
			return INCREMENT_ERROR
		default:
			for i := fst; i.GreaterThanOrEqual(lst); i = i.Add(inc) {
				f, _ := i.Float64()
				fmt.Fprintf(c.outStream, "%s%s", fmt.Sprintf(flagF, f), flagS)
			}
		}
	} else {
		f, _ := fst.Float64()
		fmt.Fprintf(c.outStream, "%s%s", fmt.Sprintf(flagF, f), flagS)
	}
	fmt.Fprint(c.outStream, flagT)
	return SUCCESS
}

func errorf(w io.Writer, format string, a ...interface{}) {
	fmt.Fprintf(w, fmt.Sprintf("goseq: %s", format), a...)
}
