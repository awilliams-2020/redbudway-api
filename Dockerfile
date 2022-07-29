FROM golang:1.18.4

COPY /etc/letsencrypt/live/redbudway.com/fullchain.pem .
COPY /etc/letsencrypt/live/redbudway.com/privkey.pem .

COPY . .
RUN go mod tidy
RUN go build cmd/redbud-way-api-server/main.go

CMD ["./main --tls-certificate=fullchain.pem --tls-key=privkey.pem"]