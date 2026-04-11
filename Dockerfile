FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /build/bin/goster-service \
    ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata curl

ENV TZ=UTC

WORKDIR /app

COPY --from=builder /build/bin/goster-service /app/goster-service
RUN chmod +x /app/goster-service

EXPOSE 8080

ENTRYPOINT ["/app/goster-service"]