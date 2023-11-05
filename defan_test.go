/*
defango - URL / IP / Email defanging with Golang. Make IoC harmless.

This repository is under MIT License https://github.com/edoardottt/defango/blob/main/LICENSE
*/

package defango_test

import (
	"net/url"
	"testing"

	"github.com/edoardottt/defango"
	"github.com/stretchr/testify/require"
)

func TestIP(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "Ipv4",
			input: "8.8.8.8",
			want:  "8[.]8[.]8[.]8",
		},
		{
			name:  "Short Ipv4",
			input: "1.1",
			want:  "1[.]1",
		},
		{
			name:  "Ipv4 + port",
			input: "8.8.8.8:53",
			want:  "8[.]8[.]8[.]8[:]53",
		},
		{
			name:  "Short Ipv4 + port",
			input: "8.8:53",
			want:  "8[.]8[:]53",
		},
		{
			name:  "Ipv6",
			input: "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			want:  "2001[:]0db8[:]85a3[:]0000[:]0000[:]8a2e[:]0370[:]7334",
		},
		{
			name:  "Short Ipv6",
			input: "2001:0db8:85a3::8a2e:0370:7334",
			want:  "2001[:]0db8[:]85a3[::]8a2e[:]0370[:]7334",
		},
		{
			name:  "Ipv6 + port",
			input: "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080",
			want:  "[2001[:]0db8[:]85a3[:]0000[:]0000[:]8a2e[:]0370[:]7334][:]8080",
		},
		{
			name:  "Short Ipv6 + port",
			input: "[2001:0db8:85a3::8a2e:0370:7334]:8080",
			want:  "[2001[:]0db8[:]85a3[::]8a2e[:]0370[:]7334][:]8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := defango.IP(tt.input)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestURL(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  string
		err   error
	}{
		{
			name:  "empty string",
			input: "",
			want:  "",
			err:   nil,
		},
		{
			name:  "https URL",
			input: "https://www.edoardoottavianelli.it",
			want:  "hxxps://www[.]edoardoottavianelli[.]it",
			err:   nil,
		},
		{
			name:  "https URL with path",
			input: "https://github.com/edoardottt/defangjs",
			want:  "hxxps://github[.]com/edoardottt/defangjs",
			err:   nil,
		},
		{
			name:  "http URL",
			input: "http://www.edoardoottavianelli.it",
			want:  "hxxp://www[.]edoardoottavianelli[.]it",
			err:   nil,
		},
		{
			name:  "http URL with path",
			input: "http://github.com/edoardottt/defangjs",
			want:  "hxxp://github[.]com/edoardottt/defangjs",
			err:   nil,
		},
		{
			name:  "http URL and port",
			input: "http://www.edoardoottavianelli.it:8080",
			want:  "hxxp://www[.]edoardoottavianelli[.]it[:]8080",
			err:   nil,
		},
		{
			name:  "http URL with path and port",
			input: "http://github.com:8000/edoardottt/defangjs",
			want:  "hxxp://github[.]com[:]8000/edoardottt/defangjs",
			err:   nil,
		},
		{
			name:  "ftp URL",
			input: "ftp://github.com/edoardottt/defangjs",
			want:  "fxp://github[.]com/edoardottt/defangjs",
			err:   nil,
		},
		{
			name:  "ftp URL and port",
			input: "ftp://github.com:8000/edoardottt/defangjs",
			want:  "fxp://github[.]com[:]8000/edoardottt/defangjs",
			err:   nil,
		},
		{
			name:  "file URL",
			input: "file:///etc/hosts",
			want:  "fxle:///etc/hosts",
			err:   nil,
		},
		{
			name:  "URL no protocol",
			input: "www.edoardoottavianelli.it",
			want:  "www[.]edoardoottavianelli[.]it",
			err:   nil,
		},
		{
			name:  "URL malformed protocol",
			input: "://www.edoardoottavianelli.it",
			want:  "://www[.]edoardoottavianelli[.]it",
			err:   nil,
		},
		{
			name:  "URL no protocol and port",
			input: "www.edoardoottavianelli.it:8080",
			want:  "www[.]edoardoottavianelli[.]it[:]8080",
			err:   nil,
		},
		{
			name:  "https URL (url.URL)",
			input: *parseURL("https://www.edoardoottavianelli.it"),
			want:  "hxxps://www[.]edoardoottavianelli[.]it",
			err:   nil,
		},
		{
			name:  "https URL with path (url.URL)",
			input: *parseURL("https://github.com/edoardottt/defangjs"),
			want:  "hxxps://github[.]com/edoardottt/defangjs",
			err:   nil,
		},
		{
			name:  "http URL (url.URL)",
			input: *parseURL("http://www.edoardoottavianelli.it"),
			want:  "hxxp://www[.]edoardoottavianelli[.]it",
			err:   nil,
		},
		{
			name:  "http URL with path (url.URL)",
			input: *parseURL("http://github.com/edoardottt/defangjs"),
			want:  "hxxp://github[.]com/edoardottt/defangjs",
			err:   nil,
		},
		{
			name:  "http URL and port (url.URL)",
			input: *parseURL("http://www.edoardoottavianelli.it:8080"),
			want:  "hxxp://www[.]edoardoottavianelli[.]it[:]8080",
			err:   nil,
		},
		{
			name:  "http URL with path and port (url.URL)",
			input: *parseURL("http://github.com:8000/edoardottt/defangjs"),
			want:  "hxxp://github[.]com[:]8000/edoardottt/defangjs",
			err:   nil,
		},
		{
			name:  "ftp URL (url.URL)",
			input: *parseURL("ftp://github.com/edoardottt/defangjs"),
			want:  "fxp://github[.]com/edoardottt/defangjs",
			err:   nil,
		},
		{
			name:  "ftp URL and port (url.URL)",
			input: *parseURL("ftp://github.com:8000/edoardottt/defangjs"),
			want:  "fxp://github[.]com[:]8000/edoardottt/defangjs",
			err:   nil,
		},
		{
			name:  "file URL (url.URL)",
			input: *parseURL("file:///etc/hosts"),
			want:  "fxle:///etc/hosts",
			err:   nil,
		},
		{
			name:  "URL no protocol (url.URL)",
			input: *parseURL("www.edoardoottavianelli.it"),
			want:  "www[.]edoardoottavianelli[.]it",
			err:   nil,
		},
		{
			name:  "URL no protocol and port (url.URL)",
			input: *parseURL("www.edoardoottavianelli.it:8080"),
			want:  "www[.]edoardoottavianelli[.]it[:]8080",
			err:   nil,
		},
		{
			name:  "float/integer input",
			input: 2,
			want:  "",
			err:   defango.ErrUnsupportedInputType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := defango.URL(tt.input)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestEmail(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "email address",
			input: "edoardott@gmail.com",
			want:  "edoardott@gmail[.]com",
		},
		{
			name:  "email link",
			input: "mailto:edoardott@gmail.com",
			want:  "mailto[:]edoardott@gmail[.]com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := defango.Email(tt.input)
			require.Equal(t, tt.want, got)
		})
	}
}

func parseURL(input string) *url.URL {
	u, _ := url.Parse(input)
	return u
}
