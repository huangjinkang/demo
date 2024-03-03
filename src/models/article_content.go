package models

// ArticleContent 映射article_content数据表的结构体
type ArticleContent struct {
	ArticleID uint64 `gorm:"primaryKey" json:"article_id"`
	Content   string `json:"content"`
}

// TableName 设置ArticleContent的表名为article_content，如果不设置默认是article_contents
func (ArticleContent) TableName() string {
	return "article_content"
}
