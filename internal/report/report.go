// Package report renders an incoming HTTP request as a whoami-style plain
// text report, additionally decoding any header value that looks like a
// compact JWT so OIDC/OAuth2 claims forwarded by a reverse proxy are easy
// to read.
package report

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/BlackDark/whoami-oidc-decoder/internal/jwtdecode"
)

// Generate renders the report for r. hostname identifies the instance that
// served the request, matching the convention of the traefik/whoami tool.
func Generate(r *http.Request, hostname string) string {
	var b strings.Builder

	fmt.Fprintf(&b, "Hostname: %s\n", hostname)
	fmt.Fprintf(&b, "Date: %s\n", time.Now().UTC().Format(time.RFC3339))
	fmt.Fprintf(&b, "Method: %s\n", r.Method)
	fmt.Fprintf(&b, "URL: %s\n", r.URL.RequestURI())
	fmt.Fprintf(&b, "Host: %s\n", r.Host)
	fmt.Fprintf(&b, "RemoteAddr: %s\n", r.RemoteAddr)

	names := make([]string, 0, len(r.Header))
	for name := range r.Header {
		names = append(names, name)
	}
	sort.Strings(names)

	b.WriteString("Headers:\n")
	for _, name := range names {
		for _, value := range r.Header[name] {
			fmt.Fprintf(&b, "  %s: %s\n", name, value)
		}
	}

	writeDecodedJWTs(&b, names, r.Header)

	return b.String()
}

func writeDecodedJWTs(b *strings.Builder, sortedNames []string, header http.Header) {
	type match struct {
		header string
		value  string
		token  string
	}

	var matches []match
	for _, name := range sortedNames {
		for _, value := range header[name] {
			for _, token := range candidateTokens(value) {
				if jwtdecode.Lookslike(token) {
					matches = append(matches, match{header: name, value: value, token: token})
				}
			}
		}
	}

	if len(matches) == 0 {
		return
	}

	b.WriteString("\nDecoded JWTs (unverified, do not trust for authorization decisions):\n")
	for _, m := range matches {
		decoded, err := jwtdecode.Decode(m.token)
		fmt.Fprintf(b, "\n  Header: %s\n", m.header)
		if err != nil {
			fmt.Fprintf(b, "    error: %v\n", err)
			continue
		}
		b.WriteString("    JOSE Header:\n")
		writeIndented(b, decoded.Header, "      ")
		b.WriteString("    Claims:\n")
		writeIndented(b, decoded.Payload, "      ")
	}
}

// candidateTokens extracts possible JWT tokens out of a header value. Most
// headers carry the token as their entire value, but some (e.g.
// Authorization: Bearer <token>) wrap it, so both forms are tried.
func candidateTokens(value string) []string {
	tokens := []string{value}
	if rest, ok := strings.CutPrefix(value, "Bearer "); ok {
		tokens = append(tokens, rest)
	}
	return tokens
}

func writeIndented(b *strings.Builder, text, indent string) {
	for _, line := range strings.Split(strings.TrimRight(text, "\n"), "\n") {
		b.WriteString(indent)
		b.WriteString(line)
		b.WriteString("\n")
	}
}
