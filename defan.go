/*
defango - URL / IP / Email defanging with Golang. Make IoC harmless.

This repository is under MIT License https://github.com/edoardottt/defango/blob/main/LICENSE
*/

package defango

import (
	"errors"
	"net/url"
	"strings"
	"unicode/utf8"
)

var (
	// ErrUnsupportedInputType is returned by URL when input is neither
	// a string, a url.URL nor a non-nil *url.URL.
	ErrUnsupportedInputType = errors.New("unsupported input type")
)

const (
	sep         = "://"
	sepToken    = "[://]"
	dotToken    = "[.]"
	colonToken  = "[:]"
	colonsToken = "[::]"

	// fullStops are the characters, besides '.', that browsers read as
	// dots in host names.
	fullStops = "\u3002\uff0e\uff61"
)

// Option configures URL and Text.
type Option func(*options)

type options struct {
	hostOnly bool
}

// HostOnly makes URL and Text defang only the schemes and authorities
// (user info, host and port) of URLs, leaving paths, queries and
// fragments as they are.
func HostOnly() Option {
	return func(o *options) { o.hostOnly = true }
}

func newOptions(opts []Option) options {
	var o options

	for _, opt := range opts {
		opt(&o)
	}

	return o
}

// IP returns a defanged IP.
// Support both IPv4 and IPv6.
func IP(input string) string {
	return defang(input)
}

// URL returns a defanged URL.
// Accept string, url.URL and *url.URL.
// The schemes http, https, ftp and file become hxxp, hxxps, fxp and fxle,
// any other scheme is followed by "[://]" instead of "://".
// The input is not validated but defanged on a best-effort basis:
// the error is non-nil only if the input type is unsupported or
// the input is a nil *url.URL.
func URL(input interface{}, opts ...Option) (string, error) {
	var s string

	switch v := input.(type) {
	case string:
		s = v
	case url.URL:
		s = v.String()
	case *url.URL:
		if v == nil {
			return "", ErrUnsupportedInputType
		}

		s = v.String()
	default:
		return "", ErrUnsupportedInputType
	}

	var b strings.Builder

	b.Grow(defangedLen(s))
	defangURL(&b, s, newOptions(opts).hostOnly)

	return b.String(), nil
}

// Email returns a defanged email address/link.
func Email(input string) string {
	return defang(input)
}

// defang returns s with its dots and colons defanged.
func defang(s string) string {
	var b strings.Builder

	b.Grow(defangedLen(s))
	defangTo(&b, s, true)

	return b.String()
}

// defangedLen returns an upper bound on the length of s once defanged.
func defangedLen(s string) int {
	return len(s) + (len(dotToken)-len("."))*(strings.Count(s, ".")+strings.Count(s, ":"))
}

// defangURL writes the URL s to b with its schemes, dots and colons
// defanged. With hostOnly, paths, queries and fragments are kept as they
// are. Text before the first scheme, or s if it has none, is defanged as
// a URL without scheme.
func defangURL(b *strings.Builder, s string, hostOnly bool) {
	start := 0

	for {
		sepStart, sepEnd, ok := nextSep(s, start)
		if !ok {
			break
		}

		scheme := schemeStart(s, start, sepStart)
		defangHierPart(b, s[start:scheme], hostOnly)

		defanged, known := defangScheme(s[scheme:sepStart])
		b.WriteString(defanged)

		if known && sepEnd-sepStart == len(sep) {
			b.WriteString(sep)
		} else {
			b.WriteString(sepToken)
		}

		start = sepEnd
	}

	defangHierPart(b, s[start:], hostOnly)
}

// defangHierPart writes s, what follows the scheme of a URL, to b.
func defangHierPart(b *strings.Builder, s string, hostOnly bool) {
	path := strings.IndexAny(s, "/?#")
	if path < 0 {
		path = len(s)
	}

	defangTo(b, s[:path], true)

	if hostOnly {
		b.WriteString(s[path:])
	} else {
		defangTo(b, s[path:], false)
	}
}

// nextSep returns the bounds of the first scheme separator in s[from:],
// either "://" or the already defanged "[://]".
func nextSep(s string, from int) (int, int, bool) {
	i := strings.Index(s[from:], sep)
	if i < 0 {
		return 0, 0, false
	}

	i += from

	if i > from && s[i-1] == '[' && strings.HasPrefix(s[i+len(sep):], "]") {
		return i - 1, i - 1 + len(sepToken), true
	}

	return i, i + len(sep), true
}

// schemeStart returns where the scheme ending at s[end] starts, not
// before lo. The scheme is empty if there is none.
func schemeStart(s string, lo, end int) int {
	i := end
	for i > lo && isSchemeChar(s[i-1]) {
		i--
	}

	// A scheme starts with a letter.
	for i < end && !isLetter(s[i]) {
		i++
	}

	return i
}

type schemePair struct {
	fanged, defanged string
}

func schemes() [4]schemePair {
	return [...]schemePair{
		{"http", "hxxp"},
		{"https", "hxxps"},
		{"ftp", "fxp"},
		{"file", "fxle"},
	}
}

// defangScheme returns the defanged form of scheme, and whether scheme
// is one that gets defanged or is already defanged.
func defangScheme(scheme string) (string, bool) {
	for _, s := range schemes() {
		if strings.EqualFold(scheme, s.fanged) {
			return s.defanged, true
		}

		if strings.EqualFold(scheme, s.defanged) {
			return scheme, true
		}
	}

	return scheme, false
}

// defangTo writes s to b with its dots and colons defanged. If s is a
// host, the other characters that browsers read as dots there are
// defanged too. Notations that are already defanged, such as "[.]" or
// "(dot)", are kept as they are.
func defangTo(b *strings.Builder, s string, host bool) {
	written := 0

	for i := 0; i < len(s); {
		n, defanged := 1, ""

		switch c := s[i]; c {
		case '.':
			defanged = dotToken
		case ':':
			if strings.HasPrefix(s[i:], "::") {
				n, defanged = len("::"), colonsToken
			} else {
				defanged = colonToken
			}
		case '[', '(', '{':
			if l, _ := token(s[i:]); l > 0 {
				n = l
			}
		default:
			if host && (c == '%' || c >= utf8.RuneSelf) {
				if l := fullStopLen(s[i:]); l > 0 {
					n, defanged = l, dotToken
				}
			}
		}

		if defanged != "" {
			b.WriteString(s[written:i])
			b.WriteString(defanged)

			written = i + n
		}

		i += n
	}

	b.WriteString(s[written:])
}

// fullStopLen returns the length of the full stop other than '.' that s
// starts with, plain or percent-encoded, if browsers read it as a dot in
// a host name, or 0.
func fullStopLen(s string) int {
	if r, n := utf8.DecodeRuneInString(s); strings.ContainsRune(fullStops, r) {
		return n
	}

	for _, stop := range [...]string{"%2e", "%e3%80%82", "%ef%bc%8e", "%ef%bd%a1"} {
		if len(s) >= len(stop) && strings.EqualFold(s[:len(stop)], stop) {
			return len(stop)
		}
	}

	return 0
}

// token returns the length of the defanged notation s starts with, such
// as "[.]", "(dot)" or "[://]", and the text it stands for.
func token(s string) (int, string) {
	if len(s) < len(dotToken) {
		return 0, ""
	}

	var closing byte

	switch s[0] {
	case '[':
		closing = ']'
	case '(':
		closing = ')'
	case '{':
		closing = '}'
	default:
		return 0, ""
	}

	// The longest notation is "[://]".
	end := 1 + strings.IndexByte(s[1:min(len(s), len(sepToken))], closing)
	if end < len("[.") {
		return 0, ""
	}

	var fanged string

	switch inner := s[1:end]; strings.ToLower(inner) {
	case ".", "dot":
		fanged = "."
	case ":", "::", sep:
		fanged = inner
	case "@", "at":
		fanged = "@"
	default:
		return 0, ""
	}

	return end + 1, fanged
}

func isLetter(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func isSchemeChar(c byte) bool {
	return isLetter(c) || isDigit(c) || c == '+' || c == '-'
}
