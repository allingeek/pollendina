FROM pollendina/debianopenssl:jessie
VOLUME ["/opt/pollendina"]
ENV COMMON_NAME=myownca
ENV KEY_NAME=mykey
ENV CERT_NAME=cacert.pem
CMD touch /opt/pollendina/index.txt && echo 01 >> /opt/pollendina/serial.txt && openssl req -x509 -config /opt/pollendina/openssl-ca.conf -newkey rsa:4096 -sha256 -nodes -out /opt/pollendina/cacert.pem -outform PEM
