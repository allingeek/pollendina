#!/bin/bash
# Required resources:
#   openssl-ca.cnf
#   cacert.pem
#   cakey.pem
#   index.txt
#   serial.txt

FILES=(openssl-ca.cnf cacert.pem cakey.pem index.txt serial.txt)
SDIR="/pollendina"
DDIR="/opt/pollendina"

for i in ${FILES[@]}; do
  # Test existence of each of the require resources
  if [ ! -f $DDIR/$i ];then
          cp $SDIR/$i $DDIR/$i
  fi
done


# openssl req -x509 -config openssl-ca.cnf -newkey rsa:4096 -sha256 -nodes -out cacert.pem -outform PEM



exec "$@" # run the default command
