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

## Configuring your reverse proxy

This tool only ever sees what's in the request headers — it can't reach back
into your identity provider. Most OIDC-auth middlewares forward a handful of
individual claims as headers (username, email, groups, ...) but **not** the
raw token itself, so put it behind a middleware that's configured to forward
the ID token (or access token) as its own header, e.g. `X-Oidc-Id-Token`.
The decoder detects any header whose value looks like a compact JWT (three
dot-separated segments), plain or `Bearer <jwt>`, so the header name doesn't
matter — just make sure one of them carries the actual token.

Example for [traefik-oidc-auth](https://github.com/sevensolutions/traefik-oidc-auth)
(Traefik plugin), forwarding claims *and* the raw ID token:

```yaml
http:
  middlewares:
    pocketid:
      plugin:
        traefik-oidc-auth:
          Secret: {{ env "POCKETID_SESSION_KEY" }}
          Provider:
            Url: {{ env "POCKETID_URL" }}
            ClientId: {{ env "POCKETID_PROXY_CLIENT_ID" }}
            ClientSecret: {{ env "POCKETID_PROXY_CLIENT_SECRET" }}
            UsePkce: true
          # request every scope whose claims you want to inspect - claims
          # not covered by a requested scope simply won't be in the token
          Scopes: ["openid", "profile", "email", "groups"]
          Headers:
            - Name: X-Oidc-Username
              Value: "{{ `{{ .claims.preferred_username }}` }}"
            - Name: X-Oidc-Email
              Value: "{{ `{{ .claims.email }}` }}"
            - Name: X-Oidc-Groups
              Value: "{{ `{{ .claims.groups }}` }}"
            # the important part: forward the raw token so the decoder can
            # show every claim in it, not just the ones mapped above
            - Name: X-Oidc-Id-Token
              Value: "{{ `{{ .idToken }}` }}"
```

Route this middleware to a `whoami-oidc-decoder` container/service the same
way you route it to your real apps, then open that route in a browser after
logging in — you'll see the individual mapped headers plus a `Decoded JWTs`
section with the full JOSE header and every claim in the token.

For other proxies (oauth2-proxy, Authentik forward-auth, etc.), the same
principle applies: configure it to forward the ID token or access token as a
plain header (or `Authorization: Bearer <token>`), then point it at this
tool.

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
