# syntax=docker/dockerfile:1

FROM golang:1.22

WORKDIR /app


COPY src/ .
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOPATH=/app go build -o /profile-app

EXPOSE 8080

CMD ["/profile-app"]
