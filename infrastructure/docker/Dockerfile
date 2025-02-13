FROM golang:1.23-bookworm AS base
RUN apt update
RUN apt install -y python3 python3-pip
RUN apt clean
WORKDIR /converter
COPY GTFS2NeTEx-converter .
RUN pip install -r requirements.txt --break-system-packages

FROM base AS dev
WORKDIR /code
ENTRYPOINT ["go", "run", ".", "/converter/GTFS2NeTEx-converter.py"]

FROM base AS build-env
WORKDIR /app
COPY src/. .
RUN go mod download
RUN go run .
RUN CGO_ENABLED=0 GOOS=linux go build -o main
