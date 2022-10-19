FROM golang:1.19.2 as build

WORKDIR /app

COPY . .

RUN go build -o RD-clone-api RD-Clone-API/cmd/api

FROM debian:bullseye-slim

RUN apt-get update -y && apt-get install ca-certificates -y \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=build /app/RD-clone-api /app/

RUN useradd -m admin
USER admin

ENTRYPOINT ["./RD-clone-api"]
