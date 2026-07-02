# whoami-oidc-decoder

A [traefik/whoami](https://github.com/traefik/whoami)-style debug HTTP server: it
echoes back the request it received, and additionally decodes any header
value that looks like a compact JWT (e.g. an OIDC ID token forwarded by an
auth middleware such as [traefik-oidc-auth](https://github.com/sevensolutions/traefik-oidc-auth)).

Decoding is unverified (no signature check) and for display only — never use
this to make authorization decisions.

## Usage

```sh
docker run -p 8080:8080 ghcr.io/blackdark/whoami-oidc-decoder:latest
curl -H "Authorization: Bearer <jwt>" http://localhost:8080/
```

`GET /healthz` returns `200` for container healthchecks.

Configuration is via environment variables:

| Variable | Default | Description |
| --- | --- | --- |
| `PORT` | `8080` | HTTP listen port |

## Development

Tooling is managed with [mise](https://mise.jdx.dev):

```sh
mise install
mise run build
mise run test
mise run lint
mise run zizmor
```

Releases are built with [goreleaser](https://goreleaser.com); pushing a
`vX.Y.Z` tag builds and publishes the multi-arch (`linux/amd64`,
`linux/arm64`) image to `ghcr.io/blackdark/whoami-oidc-decoder`.
