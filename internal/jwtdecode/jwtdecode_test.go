package jwtdecode_test

import (
	"strings"
	"testing"

	"github.com/BlackDark/whoami-oidc-decoder/internal/jwtdecode"
)

// A well-known example JWT (header: {"alg":"HS256","typ":"JWT"},
// payload: {"sub":"1234567890","name":"John Doe","iat":1516239022}).
const exampleJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
	"eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ." +
	"SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

func TestLookslike(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{"valid jwt", exampleJWT, true},
		{"empty string", "", false},
		{"plain text", "hello world", false},
		{"two segments", "abc.def", false},
		{"four segments", "abc.def.ghi.jkl", false},
		{"empty segment", "abc..def", false},
		{"invalid chars", "abc.d!f.ghi", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := jwtdecode.Lookslike(tc.in); got != tc.want {
				t.Errorf("Lookslike(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	got, err := jwtdecode.Decode(exampleJWT)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	if !strings.Contains(got.Header, `"alg": "HS256"`) {
		t.Errorf("Header = %q, want it to contain alg claim", got.Header)
	}
	if !strings.Contains(got.Payload, `"name": "John Doe"`) {
		t.Errorf("Payload = %q, want it to contain name claim", got.Payload)
	}
}

func TestDecode_NotAJWT(t *testing.T) {
	if _, err := jwtdecode.Decode("not-a-jwt"); err != jwtdecode.ErrNotAJWT {
		t.Errorf("Decode() error = %v, want ErrNotAJWT", err)
	}
}

func TestDecode_InvalidBase64(t *testing.T) {
	if _, err := jwtdecode.Decode("!!!.!!!.!!!"); err == nil {
		t.Error("Decode() error = nil, want error for invalid base64 segments")
	}
}
