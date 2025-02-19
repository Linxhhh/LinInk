package ioc

import (
	"context"
	"log"

	"github.com/Linxhhh/LinInk/search/repository/dao"
	"github.com/elastic/go-elasticsearch/v8"
)

func InitESClient() *elasticsearch.TypedClient {

	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		log.Fatal("创建 TypedClient 失败，err：", err)
	}

	ctx := context.TODO()
	ok, err := client.Ping().IsSuccess(ctx)
	if err != nil || !ok {
		log.Fatal("TypedClient Ping 失败，err：", err)
	}

	err = dao.InitES(client)
	if err != nil {
		log.Fatal("创建 ES 索引失败，err：", err)
	}
	return client
}
