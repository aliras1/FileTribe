FROM ubuntu:latest

ENV GOVERSION 1.10
ENV GOPATH /root/go

RUN export PATH=${PATH}:${GOPATH}/bin

RUN \
    apt-get update && \
    apt-get -y install git golang-1.9-go wget

RUN /usr/lib/go-1.9/bin/go get github.com/golang/glog
RUN /usr/lib/go-1.9/bin/go get github.com/whyrusleeping/tar-utils
RUN /usr/lib/go-1.9/bin/go get -u golang.org/x/crypto/...

RUN wget "https://dist.ipfs.io/go-ipfs/v0.4.14/go-ipfs_v0.4.14_linux-amd64.tar.gz" && \
    tar xvfz go-ipfs_v0.4.14_linux-amd64.tar.gz && \
    cd go-ipfs && \
    ./install.sh

ADD ./.docker/start.sh /
ADD . /root/go/src/ipfs-share
RUN cd /root/go/src/ipfs-share/main && \
    /usr/lib/go-1.9/bin/go build main.go

#ENTRYPOINT /opt/go/bin/ipfs-share

EXPOSE 3333
RUN cd /
RUN chmod 777 start.sh
ENTRYPOINT "/start.sh"