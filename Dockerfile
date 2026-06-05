# syntax=docker/dockerfile:1
# check=error=true

FROM golang:1.26.3-alpine3.23 AS base

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod,id=gomodcache \
	go mod download


FROM base AS dev

RUN --mount=type=cache,target=/go/pkg/mod,id=gomodcache \
	--mount=type=cache,target=/root/.cache/go-build,id=gobuildcache \
	<<EOF
	go install github.com/go-delve/delve/cmd/dlv@v1.26.3  
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1
EOF

COPY . .

# 8080 = api port
# 2345 = debugger port
EXPOSE 8080 2345

CMD ["sleep", "infinity"]


FROM base AS builder

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod,id=gomodcache \
	--mount=type=cache,target=/root/.cache/go-build,id=gobuildcache \
	<<EOF
    CGO_ENABLED=0 GOOS=linux \
    go build -o /app ./cmd/server
EOF


FROM alpine:3.23 AS production

RUN apk --no-cache add ca-certificates

RUN adduser -D -u 1001 nonroot
USER nonroot

COPY --from=builder /app /app

EXPOSE 8080

CMD ["/app"]
