# syntax=docker/dockerfile:1
FROM golang:1.21.0

COPY . /app
WORKDIR /app
RUN go mod download
RUN go build -o stellar-scope ./cmd/stellar-scope/main.go

CMD [ "/app/stellar-scope" ]



