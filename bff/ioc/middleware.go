package ioc

import (
	"net/http"
	"strings"
	"time"

	"github.com/Linxhhh/LinInk/bff/app/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
)

func InitMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		// 注册鉴权中间件
		middleware.AuthByJWT(),

		// 注册指标采集中间件
		(&middleware.MetricsMiddlewareBuilder{
			Namespace: "LinInk",
			Subsystem: "bff",
			Name:      "gin_http",
		}).Build(),

		// 注册链路追踪中间件
		otelgin.Middleware("LinInk"),

		// 配置 CORS
		cors.New(cors.Config{
			AllowCredentials: true,
			AllowHeaders:     []string{"Content-Type", "jwt-token"},
			ExposeHeaders:    []string{"jwt-token", "trace-id", "Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
			// AllowAllOrigins:  true,
			AllowOriginFunc: func(origin string) bool {
				// 允许开发环境的 localhost 和 127.0.0.1
				if strings.HasPrefix(origin, "http://localhost") || strings.HasPrefix(origin, "https://localhost") ||
					strings.HasPrefix(origin, "http://127.0.0.1") || strings.HasPrefix(origin, "https://127.0.0.1") {
					return true
				}
				return strings.Contains(origin, "webook.com")
			},
			MaxAge: 12 * time.Hour,
		}),

		// 处理 Options
		handleOptions(),

		// 设置 Trace Id
		setHeader(),
	}
}

// 处理 OPTIONS 请求
func handleOptions() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// 在响应头中设置 Trace Id
func setHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("trace-id", trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String())
	}
}

/* 注册 Session 会话中间件
store, err := redis.NewStore(16, "tcp", "localhost:6379", "", []byte("sgpLG7yh8mUYnh619gO0P5HdYftPKpAQ"), []byte("FlIESLxvbN5wiYZS6v7HgLkqsTmED0yh"))
if err != nil {
	panic(err)
}
router.Use(sessions.Sessions("ssid", store))
*/
