FROM debian:jessie
RUN apt-get update && apt-get upgrade -y && apt-get install -y openssl 
COPY ./service.go /pollendina/
WORKDIR /pollendina
RUN apt-get install -y golang-go && go build -v && mv pollendina /usr/bin/pollendina && apt-get remove -y golang-go
WORKDIR /opt/pollendina
CMD pollendina
# COPY ./pollendina /usr/bin/pollendina
