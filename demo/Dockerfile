FROM nginx

COPY html/*  /usr/share/nginx/html/
COPY ssl.conf /etc/nginx/conf.d/ssl.conf
COPY server.crt /certs/server.crt
COPY server.key /certs/server.key
COPY cacert.pem /certs/cacert.pem
RUN apt-get update && apt-get install -y \
	curl nano