
FROM golang:1.18-alpine

WORKDIR /usr/src/golangrms

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /docker-gs-ping

EXPOSE 8000

CMD [ "/docker-gs-ping" ]

