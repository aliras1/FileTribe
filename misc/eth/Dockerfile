FROM ubuntu:latest

ENV GOVERSION 1.10
ENV GOPATH /root/go

RUN export PATH=${PATH}:${GOPATH}/bin

RUN apt-get update
RUN apt-get install -y software-properties-common
RUN add-apt-repository ppa:ethereum/ethereum
RUN \
    apt-get update && \
    apt-get -y install git golang-1.9-go wget ethereum nodejs npm

RUN npm install -g truffle

COPY ./eth /ethcode
COPY ./misc/eth /ethnet

EXPOSE 8000
EXPOSE 30304
RUN chmod +x /ethnet/start.sh
ENTRYPOINT "/ethnet/start.sh"