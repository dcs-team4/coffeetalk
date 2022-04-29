# syntax=docker/dockerfile:1

FROM golang:1.18

ENV APP_HOME /app
WORKDIR APP_HOME

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /docker-gs-ping

EXPOSE $PORT

CMD [ "/docker-gs-ping" ]
