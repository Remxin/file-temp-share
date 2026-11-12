# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS build
WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata wget \
	&& adduser -D -u 1000 appuser \
	&& mkdir -p /data/uploads \
	&& chown -R appuser:appuser /data

WORKDIR /app
COPY --from=build /out/api /app/api

USER appuser
EXPOSE 8080
ENTRYPOINT ["/app/api"]
