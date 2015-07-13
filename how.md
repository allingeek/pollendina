# Demo: Nginx server and curl client

`cd example`

# Access example
`cd example`

# Generate root certificate and run CA. Keys will be stored in the directory
`docker run -d --name pollendina_ca -p 33004:33004 -v "$PWD":/opt/pollendina/ pollendina/pollendina`

# Distribute keys to the images

```
cp cacert.pem client/cacert.pem
cp cacert.pem server/cacert.pem
```

## Build server
`cd server`

### Build base image
`docker build -t server .`

### set up vars

```
CA_IP=192.168.59.103:33004
POLLENDINA_TOKEN=$(openssl rand -hex 32)
COMMON_NAME="hw"
```

### authorize and lunch container

```
curl --cacert "cacert.pem" --data "token=${POLLENDINA_TOKEN}&cn=${COMMON_NAME}" http://$CA_IP/v1/authorize
docker run --name nginx-server -e POLLENDINA_TOKEN="$POLLENDINA_TOKEN" -e COMMON_NAME="$COMMON_NAME" -e CA_IP="$CA_IP" -i -t server /bin/bash
```

## Build client
`cd client`

### Build base image
`docker build -t client .`

### set up vars

```
CA_IP=192.168.59.103:33004
POLLENDINA_TOKEN=$(openssl rand -hex 32)
COMMON_NAME="client"
```

### authorize and lunch container

```
curl --cacert "cacert.pem" --data "token=${POLLENDINA_TOKEN}&cn=${COMMON_NAME}" http://$CA_IP/v1/authorize
docker run -e POLLENDINA_TOKEN="$POLLENDINA_TOKEN" -e COMMON_NAME="$COMMON_NAME" -e CA_IP="$CA_IP" -i -t --link nginx-server:/hw client /bin/bash
```

## test mutual authentication from client to server

```
curl -v -s -k --key /certs/id.key --cert /certs/id.crt https://hw`
curl -v -debug -s -k --key /certs/id.key --cert /certs/id.crt -X GET https://hw
docker run -e POLLENDINA_TOKEN="$POLLENDINA_TOKEN" -e COMMON_NAME="$COMMON_NAME" -e CA_IP="$CA_IP" -i -t --link backstabbing_franklin:/hw customer /bin/bash
docker run -e POLLENDINA_TOKEN="$POLLENDINA_TOKEN" -e COMMON_NAME="$COMMON_NAME" -e CA_IP="$CA_IP" -i -t server /bin/bash
```
