package main

import (
	"demo/src/handlers"
	"demo/src/services"
	"github.com/gin-gonic/gin"
)

func SetupRouter(articleService *services.ArticleService) *gin.Engine {
	router := gin.Default()

	// api路由组 v1
	v1 := router.Group("/api/v1")
	{
		articleHandler := handlers.NewArticleHandler(articleService)
		// 新增文章
		v1.POST("/article", articleHandler.AddArticle)

		// 获取文章列表
		v1.GET("/articles", articleHandler.ListArticles)

		// 更新文章
		v1.PUT("/article/:article_id", articleHandler.UpdateArticle)
	}

	return router
}
