package ioc

import (
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	MyEtcdURL string = "localhost:12379"
)

func InitEtcdClient() *clientv3.Client {
	client, err := clientv3.NewFromURL(MyEtcdURL)
	if err != nil {
		panic(err)
	}
	return client
}