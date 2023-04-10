FROM golang:1.20

WORKDIR /opt/app

COPY . .

RUN go build cmd/main/main.go

EXPOSE 8000

CMD ["./main"]