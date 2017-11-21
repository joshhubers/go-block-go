FROM golang:1.8

WORKDIR /go/src/app

COPY . .

WORKDIR /go/src/app/src/main

RUN go-wrapper download
RUN go-wrapper install

EXPOSE 8080

CMD ["go-wrapper", "run", "echo ${HOSTIP}"]
