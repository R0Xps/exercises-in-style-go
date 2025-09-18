FROM --platform=linux/amd64 debian:stable-slim

RUN apt-get update && apt-get install -y sqlite3 libsqlite3-dev

ADD bin/* /usr/bin/
ADD examples /examples