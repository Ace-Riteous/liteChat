#!/bin/bash

reso_addr='registry.cn-hangzhou.aliyuncs.com/liteChat/user-rpc-dev'
tag='latest'

pod_ip="127.0.0.1"

comtainer_name='liteChat-user-rpc-test'

docker stop ${comtainer_name}

docker rm ${comtainer_name}

docker rmi ${reso_addr}:${tag}

docker pull ${reso_addr}:${tag}

docker run -p 10001:10000 -e POD_IP=${pod_ip} --name=${comtainer_name} -d ${reso_addr}:${tag}