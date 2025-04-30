FROM golang:1.24-alpine

WORKDIR /app

COPY . .

RUN go build -o bin/ecom cmd/main.go

CMD ["./bin/ecom"]
