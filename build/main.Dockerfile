FROM golang:1.19 AS builder


WORKDIR /opt/app

COPY . .

RUN go build cmd/main/main.go


FROM ubuntu:latest


RUN apt-get -y update && apt-get install -y tzdata



ENV dbData "postgres://docker:docker@127.0.0.1:5432/docker?pool_max_conns=10"



ENV TZ=Russina/Moscow

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone



ENV PostgresVer 14

ENV PostgresPort 5432



RUN apt-get -y update && apt-get install -y postgresql-$PostgresVer



USER postgres



RUN service postgresql start &&\
\
psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
\
createdb -O docker docker &&\
\
service postgresql stop



EXPOSE $PostgresPort



USER root



WORKDIR /usr/src/app



COPY . .

COPY --from=builder /opt/app/main .



EXPOSE 8000

ENV PGPASSWORD docker



CMD service postgresql start && psql -h localhost -d docker -U docker -p $PostgresPort -a -q -f ./build/init.sql && ./main
