FROM centurylink/openssl

MAINTAINER Pollendina <https://github.com/allingeek/pollendina>


RUN apk-install curl

COPY pollendina.sh /certs/pollendina.sh
RUN chmod +x /certs/pollendina.sh

ENTRYPOINT ["/certs/pollendina.sh"]