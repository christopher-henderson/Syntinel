#!/usr/bin/env bash

docker run --name syntinel -detach -p 80:80 -v $1:/opt/Syntinel -v /var/run/docker.sock:/var/run/docker.sock syntinel
