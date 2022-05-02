# syntax=docker/dockerfile:1

# Base image.
FROM golang:1.18

# Sets working directory.
ENV APP_HOME /app
WORKDIR APP_HOME

# Builds source files.
COPY . ./
RUN go build -o /coffeetalk-web

# Exposes port from PORT environment variable.
EXPOSE $PORT

# Run binary.
CMD [ "/coffeetalk-web" ]
