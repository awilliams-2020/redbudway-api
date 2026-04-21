FROM golang:1.24-alpine AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

ENV CGO_ENABLED=0
RUN go build -trimpath -ldflags="-s -w" -o /app ./cmd/redbud-way-api-server/main.go

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /usr/src/app

COPY --from=builder /app /usr/local/bin/app

EXPOSE 80

CMD ["/usr/local/bin/app", "--host", "0.0.0.0", "--port", "80"]
