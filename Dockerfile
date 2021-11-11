FROM golang:1.17

RUN mkdir /data

WORKDIR /data

ADD . ./

RUN go mod download

RUN go build -o main ./cmd/api

EXPOSE 8080

CMD [ "/data/main" ]