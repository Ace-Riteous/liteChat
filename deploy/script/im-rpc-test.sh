#!/bin/bash
reso_addr='registry.cn-hangzhou.aliyuncs.com/liteChat/im-rpc-dev'
tag='latest'

pod_ip="127.0.0.1"

container_name="liteChat-im-rpc-test"

docker stop ${container_name}

docker rm ${container_name}

docker rmi ${reso_addr}:${tag}

docker pull ${reso_addr}:${tag}

docker run -p 10002:10002 -e POD_IP=${pod_ip}  --name=${container_name} -d ${reso_addr}:${tag}