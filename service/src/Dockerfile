FROM pollendina/debianopenssl:jessie
ENV COUNTRY=US
ENV STATE=California
ENV CITY=SFO
ENV ORGANIZATION=None 
ENV CN=Pollendina
COPY . /pollendina/
VOLUME ["/opt/pollendina", "/var/crt", "/var/csr"]
WORKDIR /pollendina
RUN apt-get install -y golang-go && go build -v && mv pollendina /usr/bin/pollendina && apt-get remove -y golang-go
ENTRYPOINT ["./init.sh"]
CMD pollendina

