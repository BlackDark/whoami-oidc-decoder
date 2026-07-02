// Command whoami-oidc-decoder is a whoami-style debug HTTP server that
// additionally decodes any request header whose value looks like a JWT.
// It is meant to sit behind an OIDC-aware reverse proxy (e.g. Traefik) to
// make it easy to see what the identity provider is actually sending.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/BlackDark/whoami-oidc-decoder/internal/report"
)

func main() {
	addr := ":" + port()

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/", rootHandler(hostname))

	log.Printf("whoami-oidc-decoder listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil { //nolint:gosec // debug tool, timeouts add no value here
		log.Fatalf("server failed: %v", err)
	}
}

func rootHandler(hostname string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		// Reflecting request data back to the caller is the whole point of a
		// whoami-style debug tool; the response is served as text/plain, not
		// HTML, so there is no XSS risk despite the taint gosec sees here.
		if _, err := w.Write([]byte(report.Generate(r, hostname))); err != nil { //nolint:gosec // see comment above
			log.Printf("failed writing response: %v", err)
		}
	}
}

func healthzHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func port() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "8080"
}
