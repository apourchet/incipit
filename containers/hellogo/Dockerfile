FROM golang:1.7

COPY server.go .
RUN apt-get update
RUN apt-get install redis-tools -y

EXPOSE 8080

CMD go run server.go
