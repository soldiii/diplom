FROM golang:1.20.3

LABEL maintainer="Danil Malakhov"

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ./ ./

EXPOSE 8080
EXPOSE 5432

RUN apt-get update
RUN apt-get -y install postgresql-client

RUN go build -o main ./cmd/main.go

CMD ["./main"]