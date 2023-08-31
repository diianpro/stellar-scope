# syntax=docker/dockerfile:1
FROM golang:1.19-alpine

COPY . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 go build -ostellar-scope

CMD [ "/app/stellar-scope" ]



