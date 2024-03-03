package dtos

type ArticleUpdateResultData struct {
	ArticleID uint64 `json:"article_id"`
}

// ArticleUpdateResponse 响应文章更新请求的JSON数据结构体
type ArticleUpdateResponse struct {
	Data    ArticleUpdateResultData `json:"data"`
	Message string                  `json:"message"`
}
