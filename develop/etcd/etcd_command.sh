#!/bin/bash


# sudo docker run -d --restart=always --privileged=true --name etcd -p 2379:2379 -p 2380:2380 -v /opt/etcd_data:/bitnami/etcd --env ALLOW_NONE_AUTHENTICATION=yes --env ETCD_ADVERTISE_CLIENT_URLS=http://0.0.0.0:2379 --log-opt max-size=10m --log-opt max-file=1 bitnami/etcd:3.4.15
# 创建租约

etcdctl lease grant 20

#查看租约列表
etcdctl lease list

#查看信息（剩余时间）
etcdctl lease timetolive   xxxxxxx

#删除租约
etcdctl lease revoke   xxxxxxx

#保持租约始终存活
etcdctl lease keep-alive xxxxx

# 把key和租约关联
etcdctl put /user shenyi --lease=xxxxxooo

# 一旦租约过期，或被删掉,key就没了
# 可以查看该租约下的所有key
etcdctl lease timetolive   xxxxxxx --keys
