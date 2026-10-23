/*
defango - URL / IP / Email defanging with Golang. Make IoC harmless.

This repository is under MIT License https://github.com/edoardottt/defango/blob/main/LICENSE
*/

package defango_test

import (
	"testing"

	"github.com/edoardottt/defango"
	"github.com/stretchr/testify/require"
)

func TestRefang(t *testing.T) {
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
			name:  "defango output",
			input: "hxxps://evil[.]com[:]8080/a[.]exe",
			want:  "https://evil.com:8080/a.exe",
		},
		{
			name:  "other dot notations",
			input: "hXXp://evil(.)com{.}net[dot]org(DOT)io",
			want:  "http://evil.com.net.org.io",
		},
		{
			name:  "email addresses",
			input: "user[at]evil[.]com, user[@]evil(.)com",
			want:  "user@evil.com, user@evil.com",
		},
		{
			name:  "other schemes",
			input: "ldap[://]evil[.]com fxp://a[.]b FXLE:///etc",
			want:  "ldap://evil.com ftp://a.b file:///etc",
		},
		{
			name:  "Ipv6",
			input: "[2001[:]db8[::]1][:]80",
			want:  "[2001:db8::1]:80",
		},
		{
			name:  "defanged scheme name without ://",
			input: "the hxxp scheme",
			want:  "the hxxp scheme",
		},
		{
			name:  "brackets in prose",
			input: "[note] (see above) {x}",
			want:  "[note] (see above) {x}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, defango.Refang(tt.input))
		})
	}
}

func TestRefangUndoesDefang(t *testing.T) {
	for _, ioc := range []string{
		"https://evil.com:8080/a.exe?q=a:b#c",
		"ldap://evil.com/a",
		"[2001:db8::1]:8080",
		"::ffff:192.168.1.1",
		"mailto:first.last@evil.com",
		"file:///etc/hosts",
	} {
		t.Run(ioc, func(t *testing.T) {
			u, err := defango.URL(ioc)
			require.NoError(t, err)
			require.Equal(t, ioc, defango.Refang(u))
			require.Equal(t, ioc, defango.Refang(defango.Text(ioc)))
		})
	}
}
