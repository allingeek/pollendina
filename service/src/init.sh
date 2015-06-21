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

<<<<<<< HEAD
=======
 touch /opt/pollendina/index.txt

echo 01 >> /opt/pollendina/serial.txt
>>>>>>> origin/master

 openssl req -x509 -config openssl-ca.cnf -newkey rsa:4096 -sha256 -nodes -out cacert.pem -outform PEM

<<<<<<< HEAD
=======
 cp /pollendina/openssl-ca.cnf /opt/pollendina/openssl-ca.cnf
>>>>>>> origin/master


exec "$@" # run the default command
