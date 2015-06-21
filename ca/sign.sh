#!/bin/sh

#client side
KEY_NAME='server'
/usr/bin/openssl req  -new -newkey rsa:4096 -days 365 -nodes -keyout ${KEY_NAME}.key -out ${KEY_NAME}.csr

touch index.txt
echo '01' > serial.txt

openssl ca -config openssl-ca.cnf -policy signing_policy -extensions signing_req -out ${KEY_NAME}cert.pem -infiles ${KEY_NAME}.csr