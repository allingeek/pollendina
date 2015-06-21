FROM nginx

COPY html/*  /usr/share/nginx/html/
COPY ssl.conf /etc/nginx/conf.d/ssl.conf
RUN apt-get update && apt-get install -y \
	curl nano