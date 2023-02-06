# build stage
FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -v -o core /app/cmd/core/main.go
# final stage
FROM alpine:latest
EXPOSE 8080
RUN apk add -U tzdata
WORKDIR /app
COPY --from=builder /app/core .
COPY . .
ENTRYPOINT ./core