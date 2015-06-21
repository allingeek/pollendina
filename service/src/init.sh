#!/bin/bash
# Required resources:
#   openssl-ca.cnf
#   cacert.pem
#   cakey.pem
#   index.txt
#   serial.txt


# Test existence of each of the require resources

touch /opt/pollendina/index.txt

echo 01 >> /opt/pollendina/serial.txt

openssl req -x509 -config openssl-ca.cnf -newkey rsa:4096 -sha256 -nodes -out cacert.pem -outform PEM

cp /pollendina/openssl-ca.cnf /opt/pollendina/openssl-ca.cnf

exec "$@" # run the default command
