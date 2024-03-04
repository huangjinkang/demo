package main

import (
	"demo/src/repositories"
	"demo/src/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"io"
	"log"
	"os"
)

func main() {
	// 加载.env配置文件
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// 初始化数据库连接
	db := repositories.InitDB()

	// 初始化ES连接
	esClient := repositories.InitElasticsearch()

	// 初始化 Redis 连接
	rdb := repositories.InitRedis()

	// 创建服务层实例
	articleService := services.NewArticleService(db, esClient, rdb)

	// 设置日志
	f, _ := os.Create("logs/gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	gin.DefaultErrorWriter = io.MultiWriter(f, os.Stderr)

	// 使用router.go中的SetupRouter函数设置Gin路由
	router := SetupRouter(articleService)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "5001" // 默认服务端口号
	}

	// 启动服务器
	err = router.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to run server: %v\n", err)
	}
}
