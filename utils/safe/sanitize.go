package safe

import (
	"strings"
)

// Sanitize sanitizes the given string, removing elements that can be dangerous to log (e.g.: line endings).
// References:
// - https://owasp.org/www-community/attacks/Log_Injection
// - https://cwe.mitre.org/data/definitions/117.html
func Sanitize(s string) string {
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\r", "", -1)

	return s
}
