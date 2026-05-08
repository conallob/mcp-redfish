FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata

COPY mcp-redfish /usr/local/bin/mcp-redfish

ENTRYPOINT ["/usr/local/bin/mcp-redfish"]
