#!/bin/bash

reso_addr='registry.cn-hangzhou.aliyuncs.com/liteChat/user-api-dev'
tag='latest'

comtainer_name='liteChat-user-api-test'

docker stop ${comtainer_name}

docker rm ${comtainer_name}

docker rmi ${reso_addr}:${tag}

docker pull ${reso_addr}:${tag}

docker run -p 8888:8888 --name=${comtainer_name} -d ${reso_addr}:${tag}