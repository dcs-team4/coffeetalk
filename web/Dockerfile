# syntax=docker/dockerfile:1

# Base image.
FROM golang:1.18

# Sets working directory.
WORKDIR /app

# Builds source files.
COPY . ./
RUN go build -o /coffeetalk-web

# Run binary.
CMD [ "/coffeetalk-web" ]
