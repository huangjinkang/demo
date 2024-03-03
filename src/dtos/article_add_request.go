package dtos

// ArticleAddRequest 接收新增文章请求的JSON数据结构体
type ArticleAddRequest struct {
	Title   string `json:"title"`
	Picture string `json:"picture"`
	Content string `json:"content"`
}
