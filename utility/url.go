package utility

import (
	"net/url"
	"strings"
)

// RedactURL returns a copy of rawURL with userinfo credentials removed.
// If parsing fails, a best-effort string scrub is returned.
func RedactURL(rawURL string) string {
	if rawURL == "" {
		return ""
	}
	u, err := url.Parse(rawURL)
	if err != nil || u.Scheme == "" {
		return scrubUserinfo(rawURL)
	}
	if u.User != nil {
		username := u.User.Username()
		if _, hasPassword := u.User.Password(); hasPassword {
			u.User = url.UserPassword(username, "xxxxx")
		} else if username != "" {
			u.User = url.User(username)
		}
	}
	return u.String()
}

func scrubUserinfo(rawURL string) string {
	schemeIdx := strings.Index(rawURL, "://")
	if schemeIdx < 0 {
		return rawURL
	}
	rest := rawURL[schemeIdx+3:]
	at := strings.Index(rest, "@")
	if at < 0 {
		return rawURL
	}
	cred := rest[:at]
	hostPart := rest[at+1:]
	if colon := strings.Index(cred, ":"); colon >= 0 {
		return rawURL[:schemeIdx+3] + cred[:colon] + ":xxxxx@" + hostPart
	}
	return rawURL[:schemeIdx+3] + cred + "@" + hostPart
}
