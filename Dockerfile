FROM golang:1.7.1

MAINTAINER Antoine Pourchet

RUN apt-get clean && apt-get update \
    && apt-get clean -qq

RUN apt-get install -y jq
RUN apt-get clean -qq

RUN go get github.com/tools/godep
RUN go get -u github.com/kardianos/govendor
RUN ln -s /go/src/github.com/apourchet/incipit /src

# ADD . /go/src/github.com/apourchet/incipit
VOLUME ["/go/src/github.com/apourchet/incipit"]

RUN useradd -u 1000 -m ubuntu &&  echo "ubuntu ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

ENV HOME /root
ENV IN_DOCKER true
WORKDIR /go/src/github.com/apourchet/incipit
