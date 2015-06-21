#!/bin/bash
# Required resources:
#   openssl-ca.cnf
#   cacert.pem
#   cakey.pem
#   index.txt
#   serial.txt

#Files to copy in
#FILES=(openssl-ca.cnf cacert.pem cakey.pem index.txt serial.txt)
FILES=(openssl-ca.cnf index.txt serial.txt)
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

cd $DDIR
if [ ! -f $DDIR/cakey.pem ];then 
  openssl req -x509 -config openssl-ca.cnf -newkey rsa:4096 -sha256 -nodes -out cacert.pem -outform PEM
  chmod 444 cacert.pem
  chmod 400 cakey.pem
fi


CERTIFICATE_INFO="/C=$COUNTRY/ST=$STATE/L=$CITY/O=$ORGANIZATION/CN=$CN"
if [ ! -f $DDIR/servicekey.pem ];then 
  echo $CERTIFICATE_INFO
  openssl req -new -newkey rsa:4096 -days 365 -nodes -subj "${CERTIFICATE_INFO}" -keyout $DDIR/servicekey.pem -out $DDIR/servicecert.csr
  openssl ca -batch -config $DDIR/openssl-ca.cnf -policy signing_policy -extensions signing_req -out $DDIR/servicecert.pem -infiles $DDIR/servicecert.csr
fi

exec "$@" # run the default command

