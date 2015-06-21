#!/bin/bash
# Required resources:
#   openssl-ca.cnf
#   cacert.pem
#   cakey.pem
#   index.txt
#   serial.txt

#Files to copy in
FILES=(openssl-ca.cnf cacert.pem cakey.pem index.txt serial.txt)
#Source directory
SDIR="/pollendina"
#Destination directory
DDIR="/opt/pollendina"

for i in ${FILES[@]}; do
  # Test existence of each of the require resources
  if [ ! -f $DDIR/$i ];then
     cp $SDIR/$i $DDIR/$i
  fi
done

# touch /opt/pollendina/index.txt

#echo 01 >> /opt/pollendina/serial.txt

# openssl req -x509 -config openssl-ca.cnf -newkey rsa:4096 -sha256 -nodes -out cacert.pem -outform PEM

# cp /pollendina/openssl-ca.cnf /opt/pollendina/openssl-ca.cnf


exec "$@" # run the default command
