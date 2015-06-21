FROM pollendina/debianopenssl:jessie
COPY ./service.go /pollendina/
WORKDIR /pollendina
RUN apt-get install -y golang-go && go build -v && mv pollendina /usr/bin/pollendina && apt-get remove -y golang-go
WORKDIR /opt/pollendina
CMD pollendina