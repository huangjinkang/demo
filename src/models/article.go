package models

import "time"

// Article 映射article数据表的结构体
type Article struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string    `json:"title"`
	Picture   string    `json:"picture"`
	Summary   string    `json:"summary"`
	CreatedAt time.Time `gorm:"type:datetime" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime" json:"updated_at"`
}

// TableName 设置Article的表名为article，如果不设置默认是articles
func (Article) TableName() string {
	return "article"
}
