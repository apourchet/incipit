FROM golang:1.6.3

MAINTAINER Antoine Pourchet

RUN apt-get clean && apt-get update \
    && apt-get install -y apt-transport-https \
    && apt-get clean -qq

RUN apt-get install -y jq
RUN apt-get clean -qq

RUN go get github.com/tools/godep
RUN go get -u github.com/kardianos/govendor
RUN ln -s /go/src/github.com/apourchet/dummy /src

ADD . /go/src/github.com/apourchet/dummy
VOLUME ["/go/src/github.com/apourchet/dummy"]

RUN useradd -u 1000 -m ubuntu &&  echo "ubuntu ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

ENV HOME /root
ENV IN_DOCKER true
WORKDIR /go/src/github.com/apourchet/dummy
