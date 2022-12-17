FROM golang:1.19.2 as builder

WORKDIR /go/src/app

COPY . .

RUN go mod download

RUN go build -o /go/bin/app

EXPOSE 8080

CMD ["/go/bin/app"]
