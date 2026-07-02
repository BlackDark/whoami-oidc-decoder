# This Dockerfile expects prebuilt, statically-linked binaries per platform
# to already exist in the build context under $TARGETPLATFORM (see
# .goreleaser.yaml dockers_v2, which cross-compiles them and drives
# multi-arch builds via docker buildx).
FROM gcr.io/distroless/static-debian12:nonroot

ARG TARGETPLATFORM
COPY $TARGETPLATFORM/whoami-oidc-decoder /usr/local/bin/whoami-oidc-decoder

USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/whoami-oidc-decoder"]
