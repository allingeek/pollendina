# Generate root certificate and run CA
docker run -d --name pollendina_ca -p 33004:33004 -v "$PWD":/opt/pollendina/ pollendina/pollendina

# Copy the certificate to example and create pollendina client image
#copy to example

docker build -t myimage .
CA_IP=192.168.59.103:33004
POLLENDINA_TOCKEN=$(openssl rand -hex 32)
COMMON_NAME="dario"
curl --cacert "cacert.pem" --data "token=${POLLENDINA_TOCKEN}&cn=${COMMON_NAME}" http://$CA_IP/v1/authorize
docker run -e POLLENDINA_TOCKEN="$POLLENDINA_TOCKEN" -e COMMON_NAME="$COMMON_NAME" -e CA_IP="$CA_IP" -i -t myimage /bin/bash