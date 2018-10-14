FROM ubuntu:latest

ENV GOVERSION 1.10
ENV GOPATH /root/go

RUN export PATH=${PATH}:${GOPATH}/bin

RUN apt-get update
RUN apt-get install -y software-properties-common
#RUN add-apt-repository ppa:ethereum/ethereum
RUN \
    apt-get update && \
    apt-get -y install wget golang-1.9-go  net-tools netcat python3
#    apt-get -y install git golang-1.9-go wget ethereum netcat
#
#RUN /usr/lib/go-1.9/bin/go get github.com/golang/glog
#RUN /usr/lib/go-1.9/bin/go get github.com/whyrusleeping/tar-utils
#RUN /usr/lib/go-1.9/bin/go get -u golang.org/x/crypto/...
#RUN /usr/lib/go-1.9/bin/go get github.com/ethereum/go-ethereum
#RUN /usr/lib/go-1.9/bin/go get github.com/gorilla/mux
#RUN /usr/lib/go-1.9/bin/go get github.com/ipfs/go-ipfs-api
#RUN /usr/lib/go-1.9/bin/go get github.com/ugorji/go/codec

RUN wget "https://dist.ipfs.io/go-ipfs/v0.4.17/go-ipfs_v0.4.17_linux-amd64.tar.gz" && \
    tar xvfz go-ipfs_v0.4.17_linux-amd64.tar.gz && \
    cd go-ipfs && \
    ./install.sh

ADD ./.docker/start.sh /
#ADD . /root/go/src/ipfs-share

EXPOSE 3333
RUN cd /
RUN chmod +x /start.sh
ENTRYPOINT "/start.sh"

# RUN chmod +x /root/go/src/ipfs-share/.docker/start.sh
# ENTRYPOINT [ "/root/go/src/ipfs-share/.docker/start.sh" ]

# ENTRYPOINT [ "bash" ]