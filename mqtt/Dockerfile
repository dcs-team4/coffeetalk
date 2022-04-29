# syntax=docker/dockerfile:1

# Base image.
FROM golang:1.18

# Sets working directory.
ENV APP_HOME /app
WORKDIR APP_HOME

# Builds source files.
COPY . ./
RUN go build -o /coffeetalk-mqtt

EXPOSE $PORT

# Run binary.
CMD [ "/coffeetalk-mqtt" ]
