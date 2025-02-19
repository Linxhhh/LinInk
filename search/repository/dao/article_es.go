package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

const ArticleIndexName = "article_index"

type ArticleElasticDAO interface {
	InputArticle(ctx context.Context, art Article) error
	WithdrawArticle(ctx context.Context, id int64) error
	SearchArticle(ctx context.Context, expression string) ([]Article, error)
}

type articleElasticDAO struct {
	client *elastic.TypedClient
}

func NewArticleElasticDAO(client *elastic.TypedClient) ArticleElasticDAO {
	return &articleElasticDAO{
		client: client,
	}
}

func (dao *articleElasticDAO) InputArticle(ctx context.Context, art Article) error {
	_, err := dao.client.Create(ArticleIndexName, strconv.FormatInt(art.Id, 10)).
		Request(art).Do(ctx)
	return err
}

func (dao *articleElasticDAO) WithdrawArticle(ctx context.Context, id int64) error {
	_, err := dao.client.Delete(ArticleIndexName, strconv.FormatInt(id, 10)).Do(ctx)
	return err
}

func (dao *articleElasticDAO) SearchArticle(ctx context.Context, expression string) ([]Article, error) {

	// Bool 查询，组合多个条件
	boolQuery := types.NewBoolQuery()

	// 条件一：multi_match 模糊匹配 expression
	boolQuery.Must = []types.Query{
		{
			MultiMatch: &types.MultiMatchQuery{
				Query:  expression,
				Fields: []string{"title", "content"},
			},
		},
	}

	// 条件二：term 精确匹配 status = 1 的 Article
	boolQuery.Filter = []types.Query{
		{
			Term: map[string]types.TermQuery{
				"status": {Value: 1},
			},
		},
	}

	// 最终查询条件
	query := types.NewQuery()
	query.Bool = boolQuery

	resp, err := dao.client.Search().Index(ArticleIndexName).Query(query).Do(ctx)
	if err != nil {
		return nil, err
	}

	// 解析结果
	articles := make([]Article, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		var art Article
		// 将 _source 字段反序列化为 Article
		err := json.Unmarshal(hit.Source_, &art)
		if err != nil {
			return nil, fmt.Errorf("反序列化 Article 失败，err：%w", err)
		}
		articles = append(articles, art)
	}

	return articles, nil
}

type Article struct {
	Id       int64     `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Status   int32     `json:"status"`
	AuthorId int64     `json:"authorId"`
	Ctime    time.Time `json:"ctime"`
	Utime    time.Time `json:"utime"`
}
