# syntax=docker/dockerfile:1

# Base image.
FROM python:3.10-slim

# Sets working directory.
WORKDIR /app

# Adds source files (ignores files in .dockerignore).
COPY . ./

# Installs dependencies.
RUN pip install --no-cache-dir -r requirements.txt

# Runs web server.
CMD python3 main.py
