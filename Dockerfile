# syntax=docker/dockerfile:1

## Build
FROM golang:1.19.0-buster AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN make build

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /app/bin/api /main
COPY --from=build /app/.env .env

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/main"]