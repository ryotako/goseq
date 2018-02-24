package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}

	tests := []struct {
		input, want string
		err         bool
	}{
		// positive integers
		{input: "5", want: "1,2,3,4,5,"},
		{input: "2 5", want: "2,3,4,5,"},
		{input: "2 2 5", want: "2,4,"},
		{input: "0", want: ""},
		{input: "10", want: "1,2,3,4,5,6,7,8,9,10,"},
		{input: "10 10 50", want: "10,20,30,40,50,"},
		// negative integers
		{input: "-1", want: ""},
		{input: "1 -2", want: ""},
		{input: "-1 -1 -3", want: "-1,-2,-3,"},
		// floating point numbers
		{input: "0.1", want: ""},
		{input: "1.1", want: "1,"},
		{input: "-0.1 1", want: "-0.1,0.9,"},
		{input: "0 0.1 1", want: "0.0,0.1,0.2,0.3,0.4,0.5,0.6,0.7,0.8,0.9,1.0,"},
		// invalid inputs
		{input: "", err: true},
		{input: "a", err: true},
		{input: ".", err: true},
		{input: "-", err: true},
		{input: "+", err: true},
		// -s option
		{input: "-s @ 3", want: "1@2@3@"},
		{input: "-s <> 3", want: "1<>2<>3<>"},
		// -t option
		{input: "-t @ 3", want: "1,2,3,@"},
		{input: "-t done 3", want: "1,2,3,done"},
		// -f option
		{input: "-f %.1f 3", want: "1.0,2.0,3.0,"},
		{input: "-f %+0.2f -1 1", want: "-1.00,+0.00,+1.00,"},
	}

	for _, test := range tests {
		outStream.Reset()
		errStream.Reset()

		status := cli.Run(strings.Fields(test.input))

		if test.err {
			// 異常終了が期待される場合は，エラーメッセージが付されているかチェック
			if status == 0 {
				t.Errorf("%q >> status code should be non-zero", test.input)
			}
			if len(errStream.String()) == 0 {
				t.Error("%q >> error message is required", test.input)
			}

		} else {
			// 正常終了が期待される場合は，出力をチェック
			if status != 0 {

				t.Errorf("%q >> status code %d should be zero", test.input, status)
			}
			got := outStream.String()
			want := strings.Replace(test.want, ",", "\n", -1)
			if want != got {
				t.Errorf("%q >> want %q, but %q", test.input, want, got)
			}

		}
	}
}
