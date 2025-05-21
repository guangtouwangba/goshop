package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/goshop/pkg/config"
	"github.com/yourusername/goshop/pkg/logger"
	"go.uber.org/zap"
)

const serviceName = "gateway"

func main() {
	// 加载配置
	cfg, err := config.Load(serviceName, "")
	if err != nil {
		fmt.Printf("无法加载配置: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	log, err := logger.New(serviceName, cfg.Service.LogLevel)
	if err != nil {
		fmt.Printf("无法初始化日志: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	log.Info(ctx, "启动 API 网关",
		zap.String("environment", cfg.Service.Environment),
		zap.Int("port", cfg.HTTP.Port),
	)

	// 初始化 Gin 路由
	if cfg.Service.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// 设置全局中间件
	setupMiddlewares(router)

	// 注册路由
	setupRoutes(router)

	// 创建 HTTP 服务器
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler: router,
	}

	// 优雅启动服务器
	go func() {
		log.Info(ctx, "启动 HTTP 服务器", zap.Int("port", cfg.HTTP.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(ctx, "HTTP 服务器启动失败", zap.Error(err))
		}
	}()

	// 等待中断信号终止服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info(ctx, "接收到关闭信号")

	// 优雅关闭服务器
	log.Info(ctx, "正在关闭服务器")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(ctx, "服务器强制关闭", zap.Error(err))
	}

	log.Info(ctx, "服务器已成功关闭")
}

// 设置中间件
func setupMiddlewares(router *gin.Engine) {
	// 跨域设置
	router.Use(corsMiddleware())

	// 安全中间件
	router.Use(securityMiddleware())

	// 请求ID
	router.Use(requestIDMiddleware())

	// 其他中间件...
}

// CORS 中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// 安全中间件
func securityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}

// 请求ID中间件
func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 X-Request-ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// 生成新的请求ID
			requestID = fmt.Sprintf("%d", time.Now().UnixNano())
		}

		// 设置请求ID到上下文
		c.Set("RequestID", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}

// 设置路由
func setupRoutes(router *gin.Engine) {
	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})

	// API 版本路由
	v1 := router.Group("/api/v1")
	{
		// 用户服务路由
		userRoutes := v1.Group("/users")
		{
			userRoutes.POST("/register", forwardToService("user", "/api/v1/users/register"))
			userRoutes.POST("/login", forwardToService("user", "/api/v1/users/login"))
			userRoutes.POST("/reset-password", forwardToService("user", "/api/v1/users/reset-password"))
			userRoutes.GET("/me", authMiddleware(), forwardToService("user", "/api/v1/users/me"))
			userRoutes.PUT("/me", authMiddleware(), forwardToService("user", "/api/v1/users/me"))
			userRoutes.GET("/me/addresses", authMiddleware(), forwardToService("user", "/api/v1/users/me/addresses"))
			userRoutes.POST("/me/addresses", authMiddleware(), forwardToService("user", "/api/v1/users/me/addresses"))
		}

		// 商品服务路由
		productRoutes := v1.Group("/products")
		{
			productRoutes.GET("", forwardToService("product", "/api/v1/products"))
			productRoutes.GET("/:id", forwardToService("product", "/api/v1/products/:id"))
			productRoutes.GET("/categories", forwardToService("product", "/api/v1/products/categories"))
			productRoutes.GET("/search", forwardToService("product", "/api/v1/products/search"))
		}

		// 订单与购物车服务路由
		orderRoutes := v1.Group("/orders")
		{
			orderRoutes.POST("", authMiddleware(), forwardToService("order", "/api/v1/orders"))
			orderRoutes.GET("", authMiddleware(), forwardToService("order", "/api/v1/orders"))
			orderRoutes.GET("/:id", authMiddleware(), forwardToService("order", "/api/v1/orders/:id"))
		}

		cartRoutes := v1.Group("/cart")
		{
			cartRoutes.GET("", forwardToService("order", "/api/v1/cart"))
			cartRoutes.POST("/items", forwardToService("order", "/api/v1/cart/items"))
			cartRoutes.PUT("/items/:id", forwardToService("order", "/api/v1/cart/items/:id"))
			cartRoutes.DELETE("/items/:id", forwardToService("order", "/api/v1/cart/items/:id"))
		}

		// 支付服务路由
		paymentRoutes := v1.Group("/payments")
		{
			paymentRoutes.POST("", authMiddleware(), forwardToService("payment", "/api/v1/payments"))
			paymentRoutes.GET("/:id", authMiddleware(), forwardToService("payment", "/api/v1/payments/:id"))
			paymentRoutes.POST("/:id/refund", authMiddleware(), forwardToService("payment", "/api/v1/payments/:id/refund"))
		}

		// 营销服务路由
		marketingRoutes := v1.Group("/marketing")
		{
			marketingRoutes.GET("/coupons", forwardToService("marketing", "/api/v1/marketing/coupons"))
			marketingRoutes.POST("/coupons/validate", forwardToService("marketing", "/api/v1/marketing/coupons/validate"))
			marketingRoutes.GET("/promotions", forwardToService("marketing", "/api/v1/marketing/promotions"))
		}

		// 内容管理服务路由
		cmsRoutes := v1.Group("/cms")
		{
			cmsRoutes.GET("/pages/:slug", forwardToService("cms", "/api/v1/cms/pages/:slug"))
			cmsRoutes.GET("/posts", forwardToService("cms", "/api/v1/cms/posts"))
			cmsRoutes.GET("/posts/:slug", forwardToService("cms", "/api/v1/cms/posts/:slug"))
			cmsRoutes.GET("/banners", forwardToService("cms", "/api/v1/cms/banners"))
		}
	}
}

// 转发请求到对应服务
func forwardToService(service, path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 在实际实现中，这里应使用反向代理将请求转发到对应微服务
		// 这里只是一个简化的示例
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("请求已转发到 %s 服务，路径为 %s", service, path),
		})
	}
}

// 身份验证中间件
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "未提供认证令牌",
			})
			return
		}

		// 在真实实现中，这里应验证 JWT 令牌
		// 这里只是一个简化的示例
		c.Set("UserID", uint(1))
		c.Next()
	}
}
