FROM --platform=$BUILDPLATFORM golang:1.26.4-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod download

COPY cmd ./cmd
COPY internal ./internal

ARG TARGETOS=linux
ARG TARGETARCH
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    if [ -n "$TARGETARCH" ]; then export GOARCH="$TARGETARCH"; fi; \
    CGO_ENABLED=0 GOOS="$TARGETOS" go build -trimpath -ldflags="-s -w" -o /out/subscriptions-app ./cmd/app

RUN mkdir -p /out/data/uploads

FROM scratch

COPY --from=builder --chown=65532:65532 /out/data /data
COPY --from=builder /out/subscriptions-app /usr/local/bin/subscriptions-app

WORKDIR /data

EXPOSE 8800
USER 65532:65532

ENTRYPOINT ["/usr/local/bin/subscriptions-app"]
