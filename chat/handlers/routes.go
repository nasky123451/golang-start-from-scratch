package handlers

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"example.com/m/chat/middlewares"
)

func SetupRoutes(r *gin.Engine) {
	// CSRF 保護
	//r.Use(gin.WrapH(csrf.Protect([]byte("32-byte-long-auth-key"), csrf.Secure(false))(r)))

	// 添加 CORS 支持
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                             // 替换为你的前端地址
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},                        // 确保允许 OPTIONS 方法
		AllowHeaders:     []string{"Content-Type", "X-CSRF-Token", "Authorization"}, // 添加您需要的自定义头
		AllowCredentials: true,
	}))

	// 处理 OPTIONS 请求
	r.OPTIONS("/register", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*") // 允许所有源
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		c.Status(http.StatusNoContent) // 返回 204 No Content
	})

	// 路由设置
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.POST("/register", RegisterUser)
	r.POST("/login", LoginUser)
	r.POST("/logout", LogoutUser)

	r.GET("/ws", HandleWebSocket)

	// 使用 JWT 中间件保护以下路由
	protected := r.Group("/")
	protected.Use(middlewares.MiddlewareJWT())
	{
		protected.OPTIONS("/online-users", func(c *gin.Context) {
			c.Status(http.StatusNoContent)
		})

		protected.OPTIONS("/latest-chat-date", func(c *gin.Context) {
			c.Status(http.StatusNoContent)
		})

		protected.GET("/online-users", GetOnlineUsers)
		protected.GET("/chat-history", GetChatHistory)
		protected.GET("/latest-chat-date", GetLatestChatDate)
	}
}
