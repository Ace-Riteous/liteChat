#!/bin/bash
reso_addr='registry.cn-hangzhou.aliyuncs.com/liteChat/im-api-dev'
tag='latest'

container_name="liteChat-im-api-test"

docker stop ${container_name}

docker rm ${container_name}

docker rmi ${reso_addr}:${tag}

docker pull ${reso_addr}:${tag}

docker run -p 8882:8882  --name=${container_name} -d ${reso_addr}:${tag}