# Build stage
FROM golang:1.26.8-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -tags timetzdata \
    -ldflags="-w -s" \
    -trimpath \
    -o estimation-service \
    ./cmd/api

# Runtime stage
FROM alpine:3.20

WORKDIR /app

RUN adduser -D -u 1000 appuser

USER appuser

COPY --from=builder /build/estimation-service /usr/local/bin/estimation-service

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/estimation-service"]
