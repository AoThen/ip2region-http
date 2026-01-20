FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
# No go.sum needed - no external dependencies

COPY xdb ./xdb
COPY config.go ./config.go
COPY download.go ./download.go
COPY main.go ./main.go

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ip2region-http .

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/ip2region-http .

# Create data directory
RUN mkdir -p /app/data

COPY data/ip2region_n.xdb /app/data/ip2region_n.xdb

ENV SERVER_PORT=8999
ENV DB_PATH=/app/data/ip2region_n.xdb
ENV DOWNLOAD_MODE=true

EXPOSE 8999

HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
    CMD wget -q --spider http://localhost:8999/health || exit 1

CMD ["./ip2region-http"]
