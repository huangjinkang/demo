package handlers

import (
	"demo/src/dtos"
	"demo/src/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ArticleHandler struct {
	service *services.ArticleService
}

func NewArticleHandler(service *services.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		service: service,
	}
}

// AddArticle 处理新增文章请求
func (h *ArticleHandler) AddArticle(c *gin.Context) {
	var articleReq dtos.ArticleAddRequest
	if err := c.ShouldBindJSON(&articleReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 添加文章
	ctx := c.Request.Context()
	articleID, err := h.service.AddArticle(ctx, &articleReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dtos.ArticleAddResponse{
		Data: dtos.ArticleAddResultData{
			ArticleID: articleID,
		},
		Message: "Article added successfully.",
	}

	c.JSON(http.StatusOK, response)
}

// ListArticles 处理获取文章列表请求
func (h *ArticleHandler) ListArticles(c *gin.Context) {
	var req dtos.ArticleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取文章列表
	articles, err := h.service.ListArticles(c.Request.Context(), req.Page, req.PageSize, req.Sort, req.Order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, articles)
}

// UpdateArticle 处理更新文章请求
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	// 获取文章ID
	articleID := c.Param("article_id")

	// 验证文章ID
	id, err := strconv.ParseUint(articleID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	var articleReq dtos.ArticleUpdateRequest
	if err := c.ShouldBindJSON(&articleReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新文章
	ctx := c.Request.Context()
	err = h.service.UpdateArticle(ctx, id, &articleReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := dtos.ArticleUpdateResponse{
		Data: dtos.ArticleUpdateResultData{
			ArticleID: id,
		},
		Message: "Article updated successfully.",
	}

	c.JSON(http.StatusOK, response)
}
