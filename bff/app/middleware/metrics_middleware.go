package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricsMiddlewareBuilder struct {
	Namespace string
	Subsystem string
	Name      string
}

func (b *MetricsMiddlewareBuilder) Build() gin.HandlerFunc {


	// summary 统计响应时间
	s := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: b.Namespace,
		Subsystem: b.Subsystem,
		Name:      b.Name + "_resp_time",

		Objectives: map[float64]float64{
			0.50: 0.050, // 50% 的情况，误差为 5.0%
			0.75: 0.025, // 75% 的情况，误差为 2.5%
			0.90: 0.010, // 90% 的情况，误差为 1.0%
			0.95: 0.005, // 95% 的情况，误差为 0.5%
			0.99: 0.001, // 99% 的情况，误差为 0.1%
		},
	}, []string{"method", "path", "status"})

	prometheus.MustRegister(s)

	// guage 统计活跃请求数
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: b.Namespace,
		Subsystem: b.Subsystem,
		Name:      b.Name + "_active_req",
	})

	prometheus.MustRegister(g)

	return func(ctx *gin.Context) {

		before := time.Now()
		g.Inc()
		defer func ()  {
			// respTime := time.Now().Sub(before)
			respTime := time.Since(before)
			g.Dec()
			s.WithLabelValues(
				ctx.Request.Method, 
				ctx.FullPath(),
				strconv.Itoa(ctx.Writer.Status()),
			).Observe(float64(respTime.Milliseconds()))
		}()
		ctx.Next()
	}
}

/*
	其它健康指标：
	1. 业务错误码的统计与监控（需要反序列化响应体）
	2. 第三方服务的监控（短信、redis、kafka）
*/