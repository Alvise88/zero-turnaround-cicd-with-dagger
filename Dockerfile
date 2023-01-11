# syntax = docker/dockerfile:1.4

ARG GO_VERSION="1.19.5"

FROM golang:${GO_VERSION}-alpine AS build
WORKDIR /src
RUN apk add --no-cache file git
ENV GOMODCACHE /root/.cache/gocache
RUN --mount=target=. --mount=target=/root/.cache,type=cache \
    CGO_ENABLED=0 go build -o /out/calc -ldflags '-s -d -w' ./cmd/calc; \
    file /out/calc | grep "statically linked"

FROM scratch
COPY --from=build /out/calc /bin/calc
ENTRYPOINT ["/bin/calc"]