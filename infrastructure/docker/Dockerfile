FROM golang:1.23-bookworm AS dev
RUN apt update
RUN apt install -y python3 python3-pip
RUN apt-get clean
WORKDIR /converter
COPY GTFS2NeTEx-converter .
RUN pip install -r requirements.txt --break-system-packages
WORKDIR /code
ENTRYPOINT ["go", "run", ".", "/converter/GTFS2NeTEx-converter.py"]

FROM golang:1.23-bookworm AS build-env
WORKDIR /app
COPY src/. .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main

FROM python:3-alpine AS build
WORKDIR /converter
COPY GTFS2NeTEx-converter .
RUN pip install -r requirements.txt
WORKDIR /app
COPY --from=build-env /app/main .
ENTRYPOINT ["main", "/converter/GTFS2NeTEx-converter.py"]

