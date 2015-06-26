FROM pollendina/client-debian:jessie

# Copy Certificate Authority certificate to /certs
COPY cacert.pem /certs/cacert.pem

# Certificate parameters
ENV COMMON_NAME=dario COUNTRY=US STATE=California CITY=SF ORGANIZATION=Marriot COMMON_NAME=Room_Controller

RUN apt-get update && apt-get install -y \
	curl
