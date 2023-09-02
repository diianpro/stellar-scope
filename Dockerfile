# syntax=docker/dockerfile:1
FROM golang:1.21.0 as builder

COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 go build -o "stellar-scope" ./cmd/stellar-scope/main.go

FROM golang:1.21-alpine
COPY --from=builder /app/stellar-scope ./
CMD ["./stellar-scope"]



