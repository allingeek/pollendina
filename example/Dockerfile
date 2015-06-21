FROM pollendina/debian:jessie

# Certificate parameters
ENV COMMON_NAME=dario CA_IP=192.168.1.1 COUNTRY=US STATE=California CITY=SF ORGANIZATION=Marriot COMMON_NAME=Room_Controller

# Copy Certificate Authority certificate to /certs
COPY cacert.pem /certs/cacert.pem