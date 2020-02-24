FROM golang:1.13.8-alpine3.11 AS builder

WORKDIR /go/src
COPY main.go .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build \
    -a \
    -trimpath \
    -ldflags '-s -w -extldflags "-static"' \
    -tags 'netgo osusergo static_build' \
    -o /go/bin/app

# Compress the binary
FROM gruebel/upx:latest as upx
COPY --from=builder /go/bin/app /app.org
RUN upx --best --lzma -o /app /app.org

# Copying the binary into the final image
FROM scratch
COPY --from=upx /app .
CMD ["./app"]
