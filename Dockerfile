FROM golang:1.17-alpine3.15

WORKDIR app

COPY . .

RUN go build -o main main.go

EXPOSE 8080
EXPOSE 8081

CMD ["./main"]