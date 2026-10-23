/*
defango - URL / IP / Email defanging with Golang. Make IoC harmless.

This repository is under MIT License https://github.com/edoardottt/defango/blob/main/LICENSE
*/

package defango_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/edoardottt/defango"
	"github.com/stretchr/testify/require"
)

var (
	// defangedNotation matches the defanged notations, and the "://" after
	// a defanged scheme.
	defangedNotation = regexp.MustCompile(`(?i)\[(\.|:|::|://|dot|@|at)\]|\((\.|:|::|://|dot|@|at)\)|` +
		`\{(\.|:|::|://|dot|@|at)\}|(hxxps?|fxp|fxle)://`)

	// fullStop matches the ways to write a full stop that defango
	// defangs as "[.]" in host names.
	fullStop = regexp.MustCompile(`(?i)%2e|%e3%80%82|%ef%bc%8e|%ef%bd%a1|[\x{3002}\x{ff0e}\x{ff61}]`)
)

func addSeeds(f *testing.F) {
	f.Helper()

	for _, s := range []string{
		"",
		"http://malicious.example.link.com:666/m4lw4r3.exe",
		"HTTPS://user:pass@evil.com/r?u=ftp://b.com#x",
		"${jndi:ldap://evil.com/a}",
		"[2001:db8::1]:8080 ::ffff:1.2.3.4 fe80::1%eth0",
		"mailto:first.last@evil.co.uk",
		"evil\u3002com http://evil%2Ecom/%2e",
		"Visit evil.com. C2: 8.8.8.8:53, see report.pdf (http://x.zip).",
		"hxxp://evil[.]com ldap[://]x(.)y user[at]evil[dot]com",
		`{"url":"https://evil.com/x","ip":"10.0.0.1"}`,
		"::1:: ends with colons that are not part of the address",
		"0::#x::: has a path that ends with colons",
	} {
		f.Add(s)
	}
}

func FuzzURL(f *testing.F) {
	addSeeds(f)

	f.Fuzz(func(t *testing.T, s string) {
		full, err := defango.URL(s)
		require.NoError(t, err)

		hostOnly, err := defango.URL(s, defango.HostOnly())
		require.NoError(t, err)

		// No dot or colon is left fanged.
		for _, out := range []string{full, defango.IP(s), defango.Email(s)} {
			left := defangedNotation.ReplaceAllString(out, "")
			require.NotContains(t, left, ".", out)
			require.NotContains(t, left, ":", out)
		}

		// Defanging twice changes nothing.
		again, err := defango.URL(full)
		require.NoError(t, err)
		require.Equal(t, full, again)

		again, err = defango.URL(hostOnly, defango.HostOnly())
		require.NoError(t, err)
		require.Equal(t, hostOnly, again)

		require.Equal(t, defango.IP(s), defango.IP(defango.IP(s)))

		requireRefangs(t, s, full, hostOnly, defango.IP(s))
	})
}

func FuzzText(f *testing.F) {
	addSeeds(f)

	f.Fuzz(func(t *testing.T, s string) {
		full := defango.Text(s)
		hostOnly := defango.Text(s, defango.HostOnly())

		// Defanging twice changes nothing, and leaves nothing to defang.
		require.Equal(t, full, defango.Text(full))
		require.Equal(t, hostOnly, defango.Text(hostOnly, defango.HostOnly()))
		require.True(t, defango.IsDefanged(full))

		requireRefangs(t, s, full, hostOnly)
	})
}

// requireRefangs checks that refanging each defanged form of s gives s
// back, up to the case of schemes and the way full stops are written,
// unless s contains defanged notations of its own.
func requireRefangs(t *testing.T, s string, defanged ...string) {
	t.Helper()

	if defango.Refang(s) != s {
		return
	}

	want := fullStop.ReplaceAllString(s, ".")

	for _, d := range defanged {
		got := fullStop.ReplaceAllString(defango.Refang(d), ".")
		require.Truef(t, strings.EqualFold(want, got), "%q refangs to %q", d, defango.Refang(d))
	}
}
