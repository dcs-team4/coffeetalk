# syntax=docker/dockerfile:1

# Base image.
FROM python:3.10-slim

# For logging in GCP.
ENV PYTHONUNBUFFERED True

# Sets working directory.
ENV APP_HOME /app
WORKDIR $APP_HOME

# Adds source files (ignores files in .dockerignore).
COPY . ./

# Installs dependencies.
RUN pip install --no-cache-dir -r requirements.txt

# Runs web server.
CMD python3 main.py
