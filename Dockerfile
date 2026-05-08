FROM golang:1.25-alpine AS builder

WORKDIR /build

# Download dependencies first for layer caching
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build \
    -ldflags="-s -w" \
    -o mcp-redfish \
    .

# Use a minimal runtime image with CA certificates for TLS
FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /build/mcp-redfish /usr/local/bin/mcp-redfish

ENTRYPOINT ["/usr/local/bin/mcp-redfish"]
