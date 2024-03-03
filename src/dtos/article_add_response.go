package dtos

type ArticleAddResultData struct {
	ArticleID uint64 `json:"article_id"`
}

// ArticleAddResponse 响应新增文章请求的JSON数据结构体
type ArticleAddResponse struct {
	Data    ArticleAddResultData `json:"data"`
	Message string               `json:"message"`
}
