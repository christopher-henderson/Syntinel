#!/usr/bin/env bash

set -e

yum install -y python-setuptools python-devel epel-release gcc wget
yum install -y redis
easy_install pip
pip install -r compozed/requirements.txt


# The golang in yum is comically old. Since 1.6 there have been major
# improvements to the compiler backend and especially the garbage collector.
wget https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.8.linux-amd64.tar.gz
mkdir -p /opt/go/src
echo "export PATH=\$PATH:/usr/local/go/bin" >> /etc/profile
echo "export GOPATH=/opt/go" >> /etc/profile
source /etc/profile
cd distributorCap/src/LoadBalancer
go build -i -o /usr/sbin/LoadBalancer .
