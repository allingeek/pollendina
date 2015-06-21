FROM pollendina/debianopenssl:jessie
VOLUME ["/opt/pollendina"]
ENV COMMON_NAME=myownca
ENV KEY_NAME=mykey
ENV CERT_NAME=cacert.pem
WORKDIR /opt/pollendina
CMD touch index.txt && echo 01 >> serial.txt && openssl req -x509 -config openssl-ca.conf -newkey rsa:4096 -sha256 -nodes -out cacert.pem -outform PEM
