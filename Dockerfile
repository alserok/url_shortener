FROM golang:alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/url_shortener/main.go

FROM alpine:latest AS final

WORKDIR /app

COPY --from=builder /build/main /main

CMD ["./main"]