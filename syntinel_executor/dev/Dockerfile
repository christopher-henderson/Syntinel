FROM docker.io/centos

MAINTAINER Christopher Henderson

RUN yum install -y go git wget
COPY script.sh $HOME/script.sh
CMD chmod +x script.sh && ./script.sh
