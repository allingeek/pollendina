POLLENDINA_TOCKEN=$(openssl rand -hex 32)
curl --cert "cacert.pem" --data "token=${POLLENDINA_TOCKEN}&cn=${COMMON_NAME}" https://$CA_IP/v1/authorization
docker run -e POLLENDINA_TOCKEN="$POLLENDINA_TOCKEN" -e COMMON_NAME="$COMMON_NAME" -i -t dario/base /bin/bash