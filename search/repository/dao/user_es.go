package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

const UserIndexName = "user_index"

type UserElasticDAO interface {
	InputUser(ctx context.Context, user User) error
	SearchUser(ctx context.Context, expression string) ([]User, error)
}

type userElasticDAO struct {
	client *elastic.TypedClient
}

func NewUserElasticDAO(client *elastic.TypedClient) UserElasticDAO {
	return &userElasticDAO{
		client: client,
	}
}

func (dao *userElasticDAO) InputUser(ctx context.Context, user User) error {

	_, err := dao.client.Create(UserIndexName, strconv.FormatInt(user.Id, 10)).
		Document(user).Do(ctx)
	return err
}

func (dao *userElasticDAO) SearchUser(ctx context.Context, expression string) ([]User, error) {

	query := types.NewQuery()
	query.Match = map[string]types.MatchQuery{
		"nickName": {Query: expression},
	}

	resp, err := dao.client.Search().Index(UserIndexName).Query(query).Do(ctx)
	if err != nil {
		return nil, err
	}

	// 解析结果
	users := make([]User, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		var user User
		// 将 _source 字段反序列化为 Article
		err := json.Unmarshal(hit.Source_, &user)
		if err != nil {
			return nil, fmt.Errorf("反序列化 User 失败，err：%w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

type User struct {
	Id           int64  `json:"id"`
	NickName     string `json:"nickName"`
	Introduction string `json:"introduction"`
}
