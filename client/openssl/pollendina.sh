#!/bin/sh

#GLOBAL VARS REQUIRED: CA_IP, COMMON_NAME
KEY_NAME="id"

CERTIFICATE_INFO="/C=${COUNTRY}/ST=${STATE}/L=${CITY}/O=${ORGANIZATION}/CN=${COMMON_NAME}"

echo $CERTIFICATE_INFO

# Generate key and create CSR
openssl req  -new -newkey rsa:4096 -days 365 -nodes -subj "${CERTIFICATE_INFO}" -keyout /certs/${KEY_NAME}.key -out /certs/${KEY_NAME}.csr

echo Authenticating with token: $POLLENDINA_TOCKEN

# Send CSR to Certificate Authority
curl --cacert /certs/cacert.pem -X PUT -s -D status --data "$(cat /certs/${KEY_NAME}.csr)" http://$CA_IP/v1/sign/${POLLENDINA_TOCKEN} -o /certs/${KEY_NAME}.crt

STATUS=$(cat status | grep HTTP/1.1 | awk {'print $2'})


if [ '$STATUS'='100 200' ]; then
   echo Container key signed by Certificate Authority: Successfully
else
	echo Error while signing container key: $STATUS
fi

rm status

echo installing CA certificate

cp /certs/cacert.pem /usr/local/share/ca-certificates/pollendina.crt
update-ca-certificates

exec $@