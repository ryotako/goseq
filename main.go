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
	var flagF, flagW bool
	flags := map[string]string{
		"-f": "%g",
		"-s": "\n",
		"-t": "",
	}
	numArgs := []string{}

loop:
	for i := 0; i < len(args); i++ {
		if args[i] == "-w" {
			flagW = true
		} else if _, ok := flags[args[i]]; ok {
			if args[i] == "-f" {
				flagF = true
			}
			if i+1 < len(args) {
				flags[args[i]] = args[i+1]
				i++
			} else {
				errorf(c.errStream, "option requires an argument -- s")
			}
		} else {
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

	if !isValidFormat(flags["-f"]) {
		errorf(c.errStream, "invalid format string: `%s'", flags["-f"])
		return INVALID_SYNTAX
	}

	// values for option -w
	beforePoint := getBeforePoint(fst.String())
	if b := getBeforePoint(inc.String()); b > beforePoint {
		beforePoint = b
	}
	if b := getBeforePoint(lst.String()); b > beforePoint {
		beforePoint = b
	}
	afterPoint := getAfterPoint(fst.String())
	if a := getAfterPoint(inc.String()); a > afterPoint {
		afterPoint = a
	}
	// if a := getAfterPoint(lst.String()); a > afterPoint {
	// 	afterPoint = a
	// }

	if fst.LessThan(lst) {
		if inc.Sign() > 0 {
			for i := fst; i.LessThanOrEqual(lst); i = i.Add(inc) {
				if flagF {
					f, _ := i.Float64()
					fmt.Fprintf(c.outStream, "%s%s", fmt.Sprintf(flags["-f"], f), flags["-s"])
				} else if flagW {
					s := padZeroAfterPoint(i.String(), afterPoint)
					s = padZeroBeforePoint(s, beforePoint)
					fmt.Fprintf(c.outStream, "%s%s", s, flags["-s"])
				} else {
					s := padZeroAfterPoint(i.String(), afterPoint)
					fmt.Fprintf(c.outStream, "%s%s", s, flags["-s"])
				}
			}
		}
	} else if fst.GreaterThan(lst) {
		if inc.Sign() < 0 {
			for i := fst; i.GreaterThanOrEqual(lst); i = i.Add(inc) {
				if flagF {
					f, _ := i.Float64()
					fmt.Fprintf(c.outStream, "%s%s", fmt.Sprintf(flags["-f"], f), flags["-s"])
				} else if flagW {
					s := padZeroAfterPoint(i.String(), afterPoint)
					s = padZeroBeforePoint(s, beforePoint)
					fmt.Fprintf(c.outStream, "%s%s", s, flags["-s"])
				} else {
					s := padZeroAfterPoint(i.String(), afterPoint)
					fmt.Fprintf(c.outStream, "%s%s", s, flags["-s"])
				}
			}
		}
	} else {
		f, _ := fst.Float64()
		fmt.Fprintf(c.outStream, "%s%s", fmt.Sprintf(flags["-f"], f), flags["-s"])
	}
	fmt.Fprint(c.outStream, flags["-t"])
	return SUCCESS
}

func errorf(w io.Writer, format string, a ...interface{}) {
	fmt.Fprintf(w, fmt.Sprintf("goseq: %s", format), a...)
}

// valid format
// 1. contains one of %e, %E, %f, %g, %G once, or contains none of them
// 2. does not contains %X, X is any charactor except for e, E, f, g, G
func isValidFormat(s string) bool {
	r := regexp.MustCompile(`%%|%[-\+ #]*\d*\.?\d*.`)
	i := 0
	for _, format := range r.FindAllString(s, -1) {
		if format == "%%" {
			continue
		} else if strings.ContainsAny(format, "eEfgG") {
			i++
			if i > 1 {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func isExponentialNotation(d decimal.Decimal) bool {
	f, _ := d.Float64()
	return strings.Contains(fmt.Sprintf("%g", f), "e")
}

func getBeforePoint(s string) int {
	if strings.Contains(s, ".") {
		return strings.Index(s, ".")
	} else {
		return len(s)
	}
}

func getAfterPoint(s string) int {
	if strings.Contains(s, ".") {
		return len(s) - strings.Index(s, ".") - 1
	} else {
		return 0
	}
}

func padZeroBeforePoint(s string, n int) string {
	m := n - getBeforePoint(s)
	if m > 0 {
		return strings.Repeat("0", m) + s
	} else {
		return s
	}
}

func padZeroAfterPoint(s string, n int) string {
	m := n - getAfterPoint(s)
	if m > 0 {
		if !strings.Contains(s, ".") {
			s += "."
		}
		return s + strings.Repeat("0", m)
	} else {
		return s
	}
}
