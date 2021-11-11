FROM golang:1.17

WORKDIR /data

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build -o main ./cmd/api

EXPOSE 8080

CMD [ "/data/main" ]