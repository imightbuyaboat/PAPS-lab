FROM golang:1.24 AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY .env .
COPY static/ static/
COPY templates/ templates/

EXPOSE 8080

CMD ["./main"]