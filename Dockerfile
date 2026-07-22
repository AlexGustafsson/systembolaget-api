FROM --platform=${BUILDPLATFORM} golang:1.26.5 AS builder

WORKDIR /src

# Use the toolchain specified in go.mod, or newer
ENV GOTOOLCHAIN=auto

COPY go.mod go.sum .
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
  go mod download && go mod verify

COPY cmd cmd
COPY systembolaget systembolaget

ARG TARGETARCH
ARG TARGETOS
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
  GOARCH=${TARGETARCH} GOOS=${TARGETOS} CGO_ENABLED=0 go build -a -ldflags="-s -w" -o /usr/local/bin/systembolaget ./cmd/systembolaget/...
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
  GOARCH=${TARGETARCH} GOOS=${TARGETOS} CGO_ENABLED=0 go build -a -ldflags="-s -w" -o /usr/local/bin/proxy ./cmd/proxy/...

FROM scratch AS export

COPY --from=builder /usr/local/bin/systembolaget systembolaget
COPY --from=builder /usr/local/bin/proxy proxy

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/local/bin/systembolaget /usr/local/bin/systembolaget
COPY --from=builder /usr/local/bin/proxy /usr/local/bin/proxy

ENTRYPOINT ["systembolaget"]
