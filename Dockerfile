# syntax=docker/dockerfile:1
FROM golang:1.21.0 as builder

COPY . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 go build -o stellar-scope ./cmd/stellar-scope/main.go

FROM golang:1.21-alpine
WORKDIR /app
COPY --from=builder /app/stellar-scope .
COPY --from=builder /app/migration /app/migration
#RUN cd ./app && ls


CMD ["./stellar-scope"]



