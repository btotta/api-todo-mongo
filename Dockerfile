FROM golang:latest

WORKDIR /go/src/app

COPY . .

RUN go get -d -v ./...

RUN go build -o main ./cmd/api

EXPOSE 8080

CMD ["./main"]
