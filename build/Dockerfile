FROM golang:1.19

WORKDIR /opt/app

COPY . .

RUN go build cmd/main/main.go

EXPOSE 8000

CMD ["./main"]