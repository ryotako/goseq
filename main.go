package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"io"
	"os"
	"strconv"
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

func printNum(w io.Writer, n float64) {
	fmt.Fprintln(w, humanize.Ftoa(n))
}

func (c *CLI) Run(a []string) int {
	fst, inc, lst := 1.0, 1.0, 1.0

	i, cnt := 0, 0
	for i < len(a) {
		// n, err := strconv.Atoi(a[i])
		n, err := strconv.ParseFloat(a[i], 64)

		if err != nil {
			fmt.Fprintf(c.errStream, "goseq: invalid floating point argument %s\n", a[i])
			return 1
		}

		switch cnt {
		case 0:
			lst = n
			cnt++
		case 1:
			fst, lst = lst, n
			cnt++
		case 2:
			inc, lst = lst, n
			cnt++
		default:
			usage(c.errStream)
			return 1
		}

		i++
	}

	if cnt == 0 {
		usage(c.errStream)
		return 1
	}

	if lst > fst {
		if inc == 0 {
			fmt.Fprintln(c.errStream, "zero increment")
			return 1
		} else if inc < 0 {
			fmt.Fprintln(c.errStream, "needs positive increment")
			return 1
		} else {
			for i := fst; i <= lst; i += inc {
				printNum(c.outStream, i)
			}
		}
	} else {
		if inc == 0 {
			fmt.Fprintln(c.errStream, "zero increment")
			return 1
		} else if fst == lst {
			printNum(c.outStream, fst)
		} else if inc > 0 {
			fmt.Fprintln(c.errStream, "needs negative increment")
			return 1
		} else {
			for i := fst; i >= lst; i += inc {
				printNum(c.outStream, i)
			}
		}
	}
	return 0
}
