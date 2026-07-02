package report_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BlackDark/whoami-oidc-decoder/internal/report"
)

const exampleJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
	"eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ." +
	"SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

func TestGenerate_BasicFields(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.test/foo?bar=1", nil)
	req.RemoteAddr = "192.0.2.1:1234"

	out := report.Generate(req, "test-host")

	for _, want := range []string{
		"Hostname: test-host",
		"Method: GET",
		"URL: /foo?bar=1",
		"Host: example.test",
		"RemoteAddr: 192.0.2.1:1234",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n---\n%s", want, out)
		}
	}
}

func TestGenerate_ListsHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.test/", nil)
	req.Header.Set("X-Oidc-Email", "user@example.test")

	out := report.Generate(req, "test-host")

	if !strings.Contains(out, "X-Oidc-Email: user@example.test") {
		t.Errorf("output missing header line\n---\n%s", out)
	}
}

func TestGenerate_DecodesJWTHeader(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.test/", nil)
	req.Header.Set("X-Oidc-Id-Token", exampleJWT)

	out := report.Generate(req, "test-host")

	if !strings.Contains(out, "Decoded JWTs") {
		t.Fatalf("output missing decoded JWT section\n---\n%s", out)
	}
	if !strings.Contains(out, "X-Oidc-Id-Token") {
		t.Errorf("output missing header name in decoded section\n---\n%s", out)
	}
	if !strings.Contains(out, `"name": "John Doe"`) {
		t.Errorf("output missing decoded claim\n---\n%s", out)
	}
}

func TestGenerate_DecodesBearerAuthorizationHeader(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.test/", nil)
	req.Header.Set("Authorization", "Bearer "+exampleJWT)

	out := report.Generate(req, "test-host")

	if !strings.Contains(out, `"name": "John Doe"`) {
		t.Errorf("output missing decoded claim from Bearer header\n---\n%s", out)
	}
}

func TestGenerate_NoJWTNoDecodedSection(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.test/", nil)
	req.Header.Set("X-Oidc-Username", "plainuser")

	out := report.Generate(req, "test-host")

	if strings.Contains(out, "Decoded JWTs") {
		t.Errorf("output should not contain decoded JWT section\n---\n%s", out)
	}
}
