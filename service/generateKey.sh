#!/bin/sh

#client side
KEY_NAME='server'
openssl req  -subj "/C=US/ST=DARIO/L=DARIO/O=DARIO/CN=DARIO" -new -newkey rsa:4096 -days 365 -nodes -keyout ${KEY_NAME}.key -out ${KEY_NAME}.csr
openssl ca -config openssl-ca.cnf -policy signing_policy -extensions signing_req -out ${KEY_NAME}cert.pem -infiles ${KEY_NAME}.csr
