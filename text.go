/*
defango - URL / IP / Email defanging with Golang. Make IoC harmless.

This repository is under MIT License https://github.com/edoardottt/defango/blob/main/LICENSE
*/

package defango

import (
	"net/netip"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	mailto         = "mailto:"
	punycodePrefix = "xn--"
	ipv4Labels     = 4

	// maxIPv6Len is the length of the longest IPv6 address.
	maxIPv6Len = len("ffff:ffff:ffff:ffff:ffff:ffff:255.255.255.255")

	// terminators end a URL in text, besides white space.
	terminators = "\"'<>`"

	// trailing are the characters that, at the end of a URL in text, are
	// more likely punctuation than part of the URL. ']' is not one of them
	// as it ends IPv6 hosts, which would be left fanged otherwise.
	trailing = ".,:;!?)" + fullStops
)

// Text defangs the URLs, domains, IP addresses and email addresses in
// input, leaving the rest of the text as it is.
// Anything shaped like a domain is treated as one: file names such as
// report.pdf are defanged too.
func Text(input string, opts ...Option) string {
	hostOnly := newOptions(opts).hostOnly

	var b strings.Builder

	b.Grow(defangedLen(input))

	done := 0 // input[:done] is written to b

	for {
		sepStart, sepEnd, ok := nextSep(input, done)
		if !ok {
			break
		}

		// The whole word around a scheme separator is defanged as a URL,
		// so that nothing glued to it is left live.
		start := wordStart(input, done, sepStart)
		end, stop := wordEnd(input, sepEnd)

		defangHosts(&b, input[done:start], hostOnly)
		defangURL(&b, input[start:end], hostOnly)
		b.WriteString(input[end:stop])

		done = stop
	}

	defangHosts(&b, input[done:], hostOnly)

	return b.String()
}

// IsDefanged reports whether input contains nothing that Text would
// defang, not counting the paths, queries and fragments of URLs (see
// HostOnly).
func IsDefanged(input string) bool {
	return Text(input, HostOnly()) == input
}

// defangHosts writes s to b, defanging the domains, IP addresses and
// email addresses in it, along with the ports and paths that follow them.
func defangHosts(b *strings.Builder, s string, hostOnly bool) {
	done := 0

	for i := 0; i < len(s); {
		start, end, next := findHost(s, i, done)
		if end > 0 {
			b.WriteString(s[done:start])
			defangURL(b, s[start:end], hostOnly)

			done = end
		}

		i = next
	}

	b.WriteString(s[done:])
}

// findHost looks for a domain or an IP address starting at s[i]. If
// there is one, it returns the span to defang: the host extended to the
// whole email address, or to the port and path after it. In any case it
// returns where to look next. s[:done] must be left as it is.
func findHost(s string, i, done int) (start, end, next int) {
	r, size := utf8.DecodeRuneInString(s[i:])
	prev, _ := utf8.DecodeLastRuneInString(s[:i])
	atWordStart := !isLabelRune(prev) && prev != '.'
	notation, _ := token(s[i:])

	// An IPv6 address can start with "::", defanged or not.
	if atWordStart && (notation > 0 || r == ':' || r < utf8.RuneSelf && isHexDigit(byte(r))) {
		if n := ipv6Len(s[i:]); n > 0 && !startsWithLabel(s[i+n:]) {
			start, end = i, i+n
			if start > done && s[start-1] == '[' && strings.HasPrefix(s[end:], "]") {
				start, end = start-1, end+1
			}

			end, next = extendHost(s, end)

			return start, end, next
		}
	}

	// Looking inside a defanged notation could find what defanging it left
	// out, so it is skipped whole.
	if notation > 0 {
		return 0, 0, i + notation
	}

	if !atWordStart || !isLabelRune(r) {
		return 0, 0, i + size
	}

	host, run := scanHost(s[i:])
	if host == 0 {
		return 0, 0, i + max(run, size)
	}

	end = i + host

	// The host is either the local part or the domain of an email address.
	if strings.HasPrefix(s[end:], "@") {
		if domain, _ := scanHost(s[end+1:]); domain > 0 {
			end += len("@") + domain

			return emailStart(s, done, i), end, end
		}
	}

	if i > done && s[i-1] == '@' {
		return emailStart(s, done, i-1), end, end
	}

	end, next = extendHost(s, end)

	return i, end, next
}

// scanHost scans the labels separated by dots that s starts with. It
// returns the length of the longest domain or IPv4 address they form,
// or 0 if there is none, and the length of all the labels.
func scanHost(s string) (host, run int) {
	labels, numeric, unicodeSep := 0, true, false

	for i := 0; ; {
		j, ascii := i, true

		var digits, letters, others int

		for j < len(s) {
			r, n := utf8.DecodeRuneInString(s[j:])
			if !isLabelRune(r) {
				break
			}

			j += n

			switch {
			case unicode.Is(unicode.Cf, r):
				// Browsers ignore invisible formatting characters in host
				// names, so they do not count.
				continue
			case r < utf8.RuneSelf && isDigit(byte(r)):
				digits++
			case unicode.IsLetter(r) || unicode.IsMark(r):
				letters++
			default:
				others++
			}

			ascii = ascii && r < utf8.RuneSelf
		}

		// Full stops other than '.' only separate ASCII labels, which
		// keeps sentences written in Chinese or Japanese out.
		if digits+letters+others == 0 || unicodeSep && !ascii {
			break
		}

		labels++
		run = j
		numeric = numeric && letters+others == 0 && digits <= len("255")

		// Top-level domains are made of letters, or punycode.
		if labels > 1 && (letters > 1 && digits+others == 0 || isPunycode(s[i:j])) {
			host = j
		}

		var n int

		n, unicodeSep = hostSepLen(s[j:], ascii)
		if n == 0 {
			break
		}

		i = j + n
	}

	if labels == ipv4Labels && numeric {
		host = run
	}

	return host, run
}

// hostSepLen returns the length of the label separator s starts with,
// and whether it is a full stop other than '.'.
func hostSepLen(s string, afterASCII bool) (int, bool) {
	if strings.HasPrefix(s, ".") {
		return len("."), false
	}

	if n, fanged := token(s); fanged == "." {
		return n, false
	}

	if r, n := utf8.DecodeRuneInString(s); afterASCII && strings.ContainsRune(fullStops, r) {
		return n, true
	}

	return 0, false
}

// ipv6Len returns the length of the IPv6 address s starts with, or 0.
func ipv6Len(s string) int {
	var buf [maxIPv6Len]byte

	addr, colons, i := buf[:0], 0, 0

	for i < len(s) && len(addr) < maxIPv6Len {
		if c := s[i]; isHexDigit(c) || c == '.' || c == ':' {
			if c == ':' {
				colons++
			}

			addr = append(addr, c)
			i++

			continue
		}

		n, fanged := token(s[i:])
		if fanged != "." && fanged != ":" && fanged != "::" {
			break
		}

		colons += strings.Count(fanged, ":")
		addr = append(addr, fanged...)
		i += n
	}

	// Leave out the dots and colons after the address, such as a full stop
	// ending a sentence.
	for colons > 1 {
		if ip, err := netip.ParseAddr(string(addr)); err == nil && ip.Is6() {
			return i
		}

		c := s[i-1]
		if c != '.' && c != ':' {
			break
		}

		if c == ':' {
			colons--
		}

		addr, i = addr[:len(addr)-1], i-1
	}

	return 0
}

// extendHost returns where the port and path after a host ending at
// s[end] end, and where the punctuation trailing them ends.
func extendHost(s string, end int) (int, int) {
	end += portLen(s[end:])

	if end < len(s) && strings.IndexByte("/?#", s[end]) >= 0 {
		return wordEnd(s, end)
	}

	return end, end
}

// portLen returns the length of the port, colon included, that s starts
// with, or 0.
func portLen(s string) int {
	var n int

	switch {
	case strings.HasPrefix(s, colonToken):
		n = len(colonToken)
	case strings.HasPrefix(s, ":"):
		n = len(":")
	default:
		return 0
	}

	end := n
	for end < len(s) && isDigit(s[end]) {
		end++
	}

	if end == n || startsWithLabel(s[end:]) {
		return 0
	}

	return end
}

// emailStart returns where the email address whose local part ends at
// s[end] starts, "mailto:" included, not before lo.
func emailStart(s string, lo, end int) int {
	i := end
	for i > lo && isLocalPartChar(s[i-1]) {
		i--
	}

	if i-lo >= len(mailto) && strings.EqualFold(s[i-len(mailto):i], mailto) {
		i -= len(mailto)
	}

	return i
}

// wordStart returns where the word ending at s[end] starts, not before lo.
func wordStart(s string, lo, end int) int {
	for end > lo {
		r, n := utf8.DecodeLastRuneInString(s[lo:end])
		if isTerminator(r) {
			break
		}

		end -= n
	}

	return end
}

// wordEnd returns where the word going on at s[from] ends, leaving out
// trailing punctuation but nothing before from, and where that
// punctuation ends. The punctuation is left as it is: it holds no host,
// and defanging it could change where the word ends.
func wordEnd(s string, from int) (int, int) {
	stop := from

	for stop < len(s) {
		r, n := utf8.DecodeRuneInString(s[stop:])
		if isTerminator(r) {
			break
		}

		stop += n
	}

	end := stop

	for end > from {
		r, n := utf8.DecodeLastRuneInString(s[:end])
		if !strings.ContainsRune(trailing, r) {
			break
		}

		end -= n
	}

	return end, stop
}

func isPunycode(label string) bool {
	return len(label) > len(punycodePrefix) && strings.EqualFold(label[:len(punycodePrefix)], punycodePrefix)
}

func isLabelRune(r rune) bool {
	if r < utf8.RuneSelf {
		c := byte(r)

		return isLetter(c) || isDigit(c) || c == '-' || c == '_'
	}

	return unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsMark(r) || unicode.Is(unicode.Cf, r)
}

func startsWithLabel(s string) bool {
	r, _ := utf8.DecodeRuneInString(s)

	return isLabelRune(r)
}

func isTerminator(r rune) bool {
	return unicode.IsSpace(r) || strings.ContainsRune(terminators, r)
}

func isLocalPartChar(c byte) bool {
	return isLetter(c) || isDigit(c) || strings.IndexByte("._%+-", c) >= 0
}

func isHexDigit(c byte) bool {
	return isDigit(c) || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F'
}
