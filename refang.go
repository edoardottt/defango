/*
defango - URL / IP / Email defanging with Golang. Make IoC harmless.

This repository is under MIT License https://github.com/edoardottt/defango/blob/main/LICENSE
*/

package defango

import (
	"strings"
)

// Refang reverses defanging. It turns notations such as "[.]", "(.)",
// "{.}", "[dot]", "[:]", "[://]", "[@]" and "[at]" back into the text
// they stand for, and the schemes hxxp, hxxps, fxp and fxle back into
// http, https, ftp and file, ignoring case.
func Refang(input string) string {
	var b strings.Builder

	b.Grow(len(input))

	written := 0

	for i := 0; i < len(input); {
		j := strings.IndexAny(input[i:], "[({")
		if j < 0 {
			break
		}

		i += j

		n, fanged := token(input[i:])
		if n == 0 {
			i++

			continue
		}

		b.WriteString(input[written:i])
		b.WriteString(fanged)

		i += n
		written = i
	}

	b.WriteString(input[written:])

	return refangSchemes(b.String())
}

// refangSchemes returns s with the schemes followed by "://" refanged.
func refangSchemes(s string) string {
	var b strings.Builder

	written := 0

	for start := 0; ; {
		i := strings.Index(s[start:], sep)
		if i < 0 {
			break
		}

		i += start
		scheme := schemeStart(s, start, i)

		if refanged := refangScheme(s[scheme:i]); refanged != s[scheme:i] {
			if written == 0 {
				b.Grow(len(s)) // refanged schemes are as long as defanged ones
			}

			b.WriteString(s[written:scheme])
			b.WriteString(refanged)

			written = i
		}

		start = i + len(sep)
	}

	if written == 0 {
		return s
	}

	b.WriteString(s[written:])

	return b.String()
}

func refangScheme(scheme string) string {
	for _, s := range schemes() {
		if strings.EqualFold(scheme, s.defanged) {
			return s.fanged
		}
	}

	return scheme
}
