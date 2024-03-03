package dtos

import "demo/src/models"

type ArticleListPageData struct {
	Total     int64 `json:"total"`
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	TotalPage int   `json:"total_page"`
}

type ArticleListData struct {
	PageData ArticleListPageData `json:"page_data"`
	List     []models.Article    `json:"list"`
}

// ArticleListResponse 响应查询文章列表请求的JSON数据结构体
type ArticleListResponse struct {
	Data    ArticleListData `json:"data"`
	Message string          `json:"message"`
}
