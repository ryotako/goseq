package main

import (
	"fmt"
	"os"
	"strconv"
)

func usage() {
	fmt.Println("usage: seq [first [incr]] last")
	os.Exit(1)
}

func abend(s string) {
	fmt.Fprintln(os.Stderr, "goseq:", s)
	os.Exit(1)
}

func main() {
	fst, inc, lst := 1, 1, 1

	a := os.Args[1:]

	i, cnt := 0, 0
	for i < len(a) {
		n, err := strconv.Atoi(a[i])

		if err != nil {
			abend("invalid integer argument " + a[i])
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
			usage()
		}

		i++
	}

	if cnt == 0 {
		usage()
	}

	if lst > fst {
		if inc == 0 {
			abend("zero increment")
		} else if inc < 0 {
			abend("needs positive increment")
		} else {
			for i := fst; i < lst+1; i += inc {
				fmt.Println(i)
			}
		}
	} else {
		if inc == 0 {
			abend("zero increment")
		} else if fst == lst {
			fmt.Println(fst)
		} else if inc > 0 {
			abend("needs negative increment")
		} else {
			for i := fst; i > lst-1; i += inc {
				fmt.Println(i)
			}
		}
	}

}
