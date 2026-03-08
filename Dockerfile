FROM golang:1.20-alpine AS builder

WORKDIR /database

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o solune .

FROM alpine:latest  

WORKDIR /root/

COPY --from=builder /database/solune .

EXPOSE 9000

CMD ["./solune"]
