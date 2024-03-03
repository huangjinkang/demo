package dtos

// ArticleListRequest 用于接收查询文章列表请求的JSON数据结构体
type ArticleListRequest struct {
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=5,max=100"`
	Sort     string `form:"sort" binding:"omitempty,oneof=id created_at"`
	Order    string `form:"order" binding:"omitempty,oneof=asc desc"`
}
