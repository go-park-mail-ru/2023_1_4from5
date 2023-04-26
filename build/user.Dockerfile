FROM golang:1.20

WORKDIR /opt/app

COPY . .

RUN go build cmd/user/main.go

EXPOSE 8020

CMD ["./main"]