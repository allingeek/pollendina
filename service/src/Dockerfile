FROM pollendina/debianopenssl:jessie
COPY . /pollendina/
VOLUME ["/opt/pollendina", "/var/crt", "/var/csr"]
WORKDIR /pollendina
RUN apt-get install -y golang-go && go build -v && mv pollendina /usr/bin/pollendina && apt-get remove -y golang-go
ENTRYPOINT ["./init.sh"]
CMD pollendina

