package ioc

import (
	"context"
	"fmt"
	"log"
	"net"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

var (
	MyEtcdURL string = "localhost:12379"
	MyService string = "LinInk/interactionService"
	port      string = ":3337"
	ttl       int64  = 15
)

func InitEtcdClient() *clientv3.Client {
	client, err := clientv3.NewFromURL(MyEtcdURL)
	if  err != nil {
		panic(err)
	}
	return client
}

func RegisterToEtcd(client *clientv3.Client) error {
	
	// 节点管理模块
	manager, err := endpoints.NewManager(client, MyService)
	if err != nil {
		return err
	}

	// 创建一个租约
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	lease, err := client.Grant(ctx, ttl)
	if err != nil {
		return err
	}

	// 注册节点，并携带租约
	ip := getOutboundIP()
	key := fmt.Sprintf("%s/%s", MyService, ip+port)
	ep := endpoints.Endpoint{Addr: ip + port}
	err = manager.AddEndpoint(ctx, key, ep, clientv3.WithLease(lease.ID))
	if err != nil {
		return err
	}
	log.Println("注册节点：" + ip + port)

	// 开启续约
	ch, err := client.KeepAlive(context.Background(), lease.ID)
	if err != nil {
		return err
	}

	go func() {
		// 当 cancel 被调用时，会退出这个循环
		for chResp := range ch {
			log.Println("续约 resp: ", chResp.String())
		}
	}()
	return err
}

// getOutboundIP 获得内网 IP 地址
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}