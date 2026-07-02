// Package jwtdecode decodes the header and payload segments of a JWT for
// display purposes. It never verifies the signature and must not be used
// for anything security sensitive.
package jwtdecode

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

// ErrNotAJWT is returned when the input does not have the three
// dot-separated segments a compact JWT is made of.
var ErrNotAJWT = errors.New("jwtdecode: value is not a compact JWT")

// Decoded holds the pretty-printed header and payload of a JWT.
type Decoded struct {
	Header  string
	Payload string
}

// Lookslike reports whether s has the shape of a compact JWT: three
// non-empty, base64url-alphabet segments separated by dots. This is a
// cheap heuristic, not a validation of the token's contents.
func Lookslike(s string) bool {
	segments := strings.Split(s, ".")
	if len(segments) != 3 {
		return false
	}
	for _, seg := range segments {
		if seg == "" || !isBase64URLAlphabet(seg) {
			return false
		}
	}
	return true
}

func isBase64URLAlphabet(s string) bool {
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
		case r == '-' || r == '_' || r == '=':
		default:
			return false
		}
	}
	return true
}

// Decode decodes the header and payload segments of a compact JWT and
// pretty-prints them as JSON. The signature segment is intentionally
// ignored: this package never verifies tokens.
func Decode(token string) (Decoded, error) {
	segments := strings.Split(token, ".")
	if len(segments) != 3 {
		return Decoded{}, ErrNotAJWT
	}

	header, err := decodeSegment(segments[0])
	if err != nil {
		return Decoded{}, err
	}
	payload, err := decodeSegment(segments[1])
	if err != nil {
		return Decoded{}, err
	}

	return Decoded{Header: header, Payload: payload}, nil
}

func decodeSegment(segment string) (string, error) {
	raw, err := base64.RawURLEncoding.DecodeString(segment)
	if err != nil {
		raw, err = base64.URLEncoding.DecodeString(segment)
		if err != nil {
			return "", err
		}
	}

	var buf bytes.Buffer
	if err := json.Indent(&buf, raw, "", "  "); err != nil {
		// Not valid JSON: surface the raw decoded bytes so the caller
		// still gets something useful instead of an error.
		return string(raw), nil
	}
	return buf.String(), nil
}
