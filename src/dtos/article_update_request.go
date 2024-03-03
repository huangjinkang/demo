package dtos

// ArticleUpdateRequest 接收文章更新请求的JSON数据结构体
type ArticleUpdateRequest struct {
	Title   string `json:"title" binding:"required"` // 文章标题不能为空值
	Picture string `json:"picture"`
	Content string `json:"content"`
}
