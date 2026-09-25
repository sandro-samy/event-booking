FROM golang:1.27.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mode download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server

# -------------------------------------------------------

FROM alpine:latest

WORKDIR /app

COPY --from=builder app/server .

EXPOSE 8080

CMD ["./server"]