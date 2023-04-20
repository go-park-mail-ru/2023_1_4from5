FROM golang:1.20

WORKDIR /opt/app

COPY . .

RUN go build cmd/auth/main.go

EXPOSE 8010

CMD ["./main"]