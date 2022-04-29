# syntax=docker/dockerfile:1

# Base image.
FROM golang:1.18

# Sets working directory.
ENV APP_HOME /app
WORKDIR APP_HOME

# Downloads dependencies (separate COPY for caching).
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Builds source files.
COPY ["*.go", "*.json"] ./
RUN go build -o /coffeetalk-mqtt

EXPOSE $PORT

# Run binary.
CMD [ "/coffeetalk-mqtt" ]
