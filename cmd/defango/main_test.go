/*
defango - URL / IP / Email defanging with Golang. Make IoC harmless.

This repository is under MIT License https://github.com/edoardottt/defango/blob/main/LICENSE
*/

package main

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/require"
)

var errIO = errors.New("i/o error")

type errWriter struct{}

func (errWriter) Write([]byte) (int, error) {
	return 0, errIO
}

func TestRun(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		refang bool
		want   string
	}{
		{
			name:  "empty input",
			input: "",
			want:  "",
		},
		{
			name:  "defang, no final newline",
			input: "C2 at http://evil.com:8080\nand 10.0.0.1",
			want:  "C2 at hxxp://evil[.]com[:]8080\nand 10[.]0[.]0[.]1",
		},
		{
			name:   "refang, CRLF line endings",
			input:  "hxxp://evil[.]com\r\nuser[at]evil[.]com\r\n",
			refang: true,
			want:   "http://evil.com\r\nuser@evil.com\r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out strings.Builder

			require.NoError(t, run(strings.NewReader(tt.input), &out, tt.refang))
			require.Equal(t, tt.want, out.String())
		})
	}
}

func TestRunErrors(t *testing.T) {
	require.ErrorIs(t, run(iotest.ErrReader(errIO), io.Discard, false), errIO)
	require.ErrorIs(t, run(strings.NewReader("evil.com"), errWriter{}, false), errIO)
}
