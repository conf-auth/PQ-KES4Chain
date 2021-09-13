
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

FROM ubuntu:18.04

RUN apt-get update && apt-get install -y git cmake ninja-build libssl-dev python3-pytest python3-pytest-xdist 
RUN apt-get update && apt-get install -y cmake gcc
RUN apt-get update && apt-get install -y unzip xsltproc doxygen graphviz 
RUN apt-get update && apt-get install -y pkg-config 
RUN apt-get update && apt-get install -y vim
RUN git clone -b main https://github.com/open-quantum-safe/liboqs.git
RUN cd liboqs && mkdir build && cd build && cmake -GNinja .. -DBUILD_SHARED_LIBS=ON && ninja && ninja install
RUN sed -i  '$a\/usr/local/lib' /etc/ld.so.conf && ldconfig

# install Go env
COPY go1.15.7.linux-amd64.tar.gz /
RUN tar -C /usr/local -xzf go1.15.7.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin 
ENV GOROOT=/usr/local/go        
ENV GOPATH=/root/go
ENV PATH=$GOPATH/bin/:$PATH
RUN go version

RUN mkdir -p /chaincode/output /chaincode/input
RUN addgroup --gid 500 chaincode && useradd -d /home/chaincode -g chaincode chaincode
RUN chown -R  chaincode:chaincode /chaincode
USER chaincode
USER root
COPY payload/pkgconfig  /usr/local/lib/pkgconfig/
ENV PKG_CONFIG_PATH /usr/local/lib/pkgconfig/
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct

