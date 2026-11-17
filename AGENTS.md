# whoami-oidc-decoder

Small Go service that decodes OIDC claims for debugging reverse-proxy authentication.

## Verify

- `mise run check`: format, vet, build, golangci-lint, and race-enabled tests.
- `mise run snapshot` after changing GoReleaser or release configuration.
- `mise run zizmor` after changing workflows.

CI runs the same checks plus actionlint and strict Renovate validation. The required gate is `check`.

## Stack

Follows https://github.com/BlackDark/agent-presets/blob/main/stack/core.md.

Deviations:
- Releases keep GoReleaser's binary and image pipeline instead of migrating to release-please.
- The GoReleaser snapshot remains a separate CI job because `mise run check` covers code only.

## Git

- Conventional Commits, one topic per PR.
- Squash unless commits need separate review.
- Force-push only with `--force-with-lease` on own branches.

## Boundaries

- Do not change release tags or publish images from a local session.
- Tools come from `mise.toml`; no global installs.
