FROM golang:alpine AS base
## ---------- ARGS
ARG TARGET_API
ARG API_PORT

## ---------- ENVS
ENV TARGET_API=${TARGET_API}
ENV API_PORT=${API_PORT}

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

## ---------- BUILD
## Stage 1 - Build the binary
FROM base AS builder
WORKDIR /build
COPY . .
RUN go build cmd/${TARGET_API}/main.go && \
    chmod +x main

# Stage 2: Compress the binary using UPX
FROM alpine AS upx
RUN apk add --no-cache upx
COPY --from=builder /build/main /upx/main
RUN upx --best --lzma /upx/main -o /upx/main_compressed

## ---------- MAIN
FROM scratch AS main
WORKDIR /app

COPY --from=upx /upx/main_compressed /app/main
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/.env /app/.env

ENTRYPOINT [ "./main" ]
EXPOSE ${API_PORT}
