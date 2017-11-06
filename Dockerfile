FROM golang:1.8

WORKDIR /go/src/app

COPY . .

WORKDIR /go/src/app/src/main

RUN go-wrapper download
RUN go-wrapper install

EXPOSE 8080

WORKDIR /go/src/app

CMD ["go-wrapper", "run"]
