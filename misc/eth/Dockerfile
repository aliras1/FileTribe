FROM ubuntu:latest

ENV GOVERSION 1.10
ENV GOPATH /root/go

RUN export PATH=${PATH}:${GOPATH}/bin

RUN apt-get update
RUN apt-get install -y software-properties-common
RUN add-apt-repository ppa:ethereum/ethereum
RUN \
    apt-get update && \
    apt-get -y install git golang-1.9-go wget ethereum npm

RUN npm install -g truffle
RUN npm install -g scrypt

#RUN cd /ethcode
#RUN npm rebuild scrypt --update-binary

EXPOSE 8000
EXPOSE 8001
EXPOSE 30304

ENTRYPOINT "/ethnet/start.sh"