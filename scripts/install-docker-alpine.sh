#!/bin/bash
# =======
# version: 1
# author: sdeletang
# description: post-install alpine
# =======

# TOOLS
apk add lsof

# DOCKER 
apk add docker

addgroup username docker

rc-update add docker boot

service docker start

# DOCKER COMPOSE
apk add docker-compose

apk add py-pip python3-dev libffi-dev openssl-dev gcc libc-dev make

pip3 install docker-compose