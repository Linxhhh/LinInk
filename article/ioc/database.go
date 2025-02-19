package ioc

import (
	"github.com/Linxhhh/LinInk/article/repository/dao"
	"github.com/Linxhhh/LinInk/pkg/callbacks"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "gorm.io/plugin/opentelemetry/tracing"
	"gorm.io/plugin/prometheus"
)

func InitDB() (master *gorm.DB, slaves []*gorm.DB) {
	master, err := gorm.Open(mysql.Open("root:123456@tcp(localhost:13306)/webook"))
	if err != nil {
		panic(err)
	}

	s1, err := gorm.Open(mysql.Open("root:123456@tcp(localhost:23306)/webook"))
	if err != nil {
		panic(err)
	}

	slaves = append(slaves, s1)

	err = master.AutoMigrate(
		&dao.Article{},
		&dao.PublishedArticle{},
	)
	if err != nil {
		panic(err)
	}

	err = initPrometheus(master, s1)
	if err != nil {
		panic(err)
	}

	err = metricsSpendTime(master, s1)
	if err != nil {
		panic(err)
	}

	/*
	err = otelTracing(master, s1)
	if err != nil {
		panic(err)
	}
	*/

	return
}

func initPrometheus(dbs ...*gorm.DB) error {

	// 采集数据库指标
	for _, db := range dbs {
		err := db.Use(prometheus.New(prometheus.Config{
			DBName:          "webook",
			RefreshInterval: 15,    // 15s 采集一次
			StartServer:     false, // 不用通过 http 服务来采集
			MetricsCollector: []prometheus.MetricsCollector{
				&prometheus.MySQL{},
			},
		}))
		return err
	}
	return nil
}

// metricsSpendTime 采集数据库操作耗时
func metricsSpendTime(dbs ...*gorm.DB) error {

	cb := callbacks.NewCallbacks("article")

	// 监控数据库操作的耗时
	for _, db := range dbs {

		err := db.Callback().Create().Before("*").
			Register("prometheus_create_before", cb.Before("create"))
		if err != nil {
			return err
		}
		err = db.Callback().Create().After("*").
			Register("prometheus_create_after", cb.After("create"))
		if err != nil {
			return err
		}

		err = db.Callback().Delete().Before("*").
			Register("prometheus_delete_before", cb.Before("delete"))
		if err != nil {
			return err
		}
		err = db.Callback().Delete().After("*").
			Register("prometheus_delete_after", cb.After("delete"))
		if err != nil {
			return err
		}

		err = db.Callback().Update().Before("*").
			Register("prometheus_update_before", cb.Before("update"))
		if err != nil {
			return err
		}
		err = db.Callback().Update().After("*").
			Register("prometheus_update_after", cb.After("update"))
		if err != nil {
			return err
		}

		err = db.Callback().Query().Before("*").
			Register("prometheus_query_before", cb.Before("query"))
		if err != nil {
			return err
		}
		err = db.Callback().Query().After("*").
			Register("prometheus_query_after", cb.After("query"))
		if err != nil {
			return err
		}

		err = db.Callback().Raw().Before("*").
			Register("prometheus_raw_before", cb.Before("raw"))
		if err != nil {
			return err
		}
		err = db.Callback().Raw().After("*").
			Register("prometheus_raw_after", cb.After("raw"))
		if err != nil {
			return err
		}

		err = db.Callback().Row().Before("*").
			Register("prometheus_row_before", cb.Before("row"))
		if err != nil {
			return err
		}
		err = db.Callback().Row().After("*").
			Register("prometheus_row_after", cb.After("row"))
		if err != nil {
			return err
		}
	}
	return nil
}

/*
// otelTracing 使用 otel 追踪数据库操作
func otelTracing(dbs ...*gorm.DB) error {
	for _, db := range dbs {
		err := db.Use(tracing.NewPlugin(
			tracing.WithDBName("webook"),
			tracing.WithoutMetrics(),        // 这里不用 metrics
			tracing.WithoutQueryVariables(), // 不要记录查询参数
		))
		return err
	}
	return nil
}
*/