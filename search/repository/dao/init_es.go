package dao

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	elastic "github.com/elastic/go-elasticsearch/v8"
	"golang.org/x/sync/errgroup"
)


// InitES 创建索引
func InitES(client *elastic.TypedClient) error {

	const timeout = time.Second * 10
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var eg errgroup.Group
	eg.Go(func() error {
		return tryCreateIndex(ctx, client, UserIndexName)
	})
	eg.Go(func() error {
		return tryCreateIndex(ctx, client, ArticleIndexName)
	})

	return eg.Wait()
}

func tryCreateIndex(ctx context.Context, client *elastic.TypedClient, idxName string) error {

	exist, err := client.Indices.Exists(idxName).Do(ctx)
	if err != nil {
		return fmt.Errorf("检测 %s 是否存在失败 %w", idxName, err)
	}
	if exist {
		return nil
	}

	_, err = client.Indices.Create(idxName).Do(ctx)
	if err != nil {
		return fmt.Errorf("初始化 %s 失败 %w", idxName, err)
	}
	return err
}
