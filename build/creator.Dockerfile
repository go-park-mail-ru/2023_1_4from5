FROM golang:1.20

WORKDIR /opt/app

COPY go.* .

RUN go mod download

COPY . .

RUN go build cmd/creator/main.go

EXPOSE 8030

CMD ["./main"]