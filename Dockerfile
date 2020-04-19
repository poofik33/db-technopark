FROM golang:1.14-stretch AS build

ADD ./ /opt/build/golang

WORKDIR /opt/build/golang

RUN go install ./cmd/app

FROM ubuntu:18.04 AS release

MAINTAINER Aleksandr Postnikov


ENV PGVER 10
RUN apt -y update && apt install -y postgresql-$PGVER

USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER forums_user WITH SUPERUSER PASSWORD 'difficult_password';" &&\
    createdb -O forums_user forums &&\
    /etc/init.d/postgresql stop

RUN echo "host all all 0.0.0.0/0 md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf

RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

EXPOSE 5432

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

EXPOSE 5000

COPY --from=build go/bin/app /usr/bin/
COPY --from=build /opt/build/golang/database /database/

CMD service postgresql start && app