FROM golang:1.23.0 AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o stresstest .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/stresstest .

CMD ["./stresstest"]
ENTRYPOINT ["./stresstest"]