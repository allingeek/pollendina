FROM debian:jessie
RUN apt-get update && apt-get upgrade -y && apt-get install -y openssl
WORKDIR /opt/pollendina
CMD pollendina
COPY ./pollendina /usr/bin/pollendina
COPY ./openssl-ca.cnf /opt/pollendina/openssl-ca.conf
