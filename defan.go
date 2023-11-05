/*
defango - URL / IP / Email defanging with Golang. Make IoC harmless.

This repository is under MIT License https://github.com/edoardottt/defango/blob/main/LICENSE
*/

package defango

import (
	"errors"
	"net/url"
	"strings"
)

var (
	ErrUnsupportedInputType = errors.New("unsupported input type")
)

// IP returns a defanged IP.
// Support both IPv4 and IPv6.
func IP(input string) string {
	result := ""

	if strings.Contains(input, "::") { // IPv6 with double column
		index := strings.Index(input, "::")
		result = strings.ReplaceAll(input[:index], ":", "[:]") +
			"[::]" +
			strings.ReplaceAll(input[index+2:], ":", "[:]")
	} else { // Full IPv6 or IPv4
		result = strings.ReplaceAll(strings.ReplaceAll(input, ".", "[.]"), ":", "[:]")
	}

	return result
}

// URL returns a defanged URL.
// Accept string and url.URL.
// For security reasons, if the url is not formatted in a proper way
// the result is an empty string with a non-nil error.
func URL(input interface{}) (string, error) {
	switch v := input.(type) {
	case string:
		return defangURL(v)
	case url.URL:
		return defangURL(v.String())
	default:
		return "", ErrUnsupportedInputType
	}
}

func defangURL(input string) (string, error) {
	result := ""

	if strings.Contains(input, "://") { // handle protocol
		result += defangProtocols(input) + "://"
		result += strings.ReplaceAll(
			strings.ReplaceAll(
				input[strings.Index(input, "://")+3:], ".", "[.]"),
			":", "[:]")
	} else {
		result += strings.ReplaceAll(strings.ReplaceAll(input, ".", "[.]"), ":", "[:]")
	}

	return result, nil
}

func defangProtocols(input string) string {
	protoMap := map[string]string{
		"http":  "hxxp",
		"https": "hxxps",
		"ftp":   "fxp",
		"file":  "fxle",
	}
	proto := input[:strings.Index(input, "://")]

	protoDefanged, ok := protoMap[proto]
	if ok {
		return protoDefanged
	} else {
		return proto
	}
}

// Email returns a defanged email address/link.
func Email(input string) string {
	return strings.ReplaceAll(strings.ReplaceAll(input, ".", "[.]"), ":", "[:]")
}
