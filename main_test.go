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
		{input: "", want: "", err: true},
		{input: "5", want: "1\n2\n3\n4\n5\n", err: false},
		{input: "2 5", want: "2\n3\n4\n5\n", err: false},
		{input: "2 2 5", want: "2\n4\n", err: false},
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
			if test.want != got {
				t.Errorf("%q >> want %q, but %q", test.input, test.want, got)
			}

		}
	}
}
