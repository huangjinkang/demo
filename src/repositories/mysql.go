package repositories

import (
	"demo/src/models"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

type MySQLRepository struct {
	db *gorm.DB
}

func NewMySQLRepository(db *gorm.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

// InitDB 初始化数据库连接
func InitDB() *gorm.DB {
	// 从环境变量中获取数据库连接配置
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	// 构建连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)

	// 使用 gorm 打开数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Database connected")
	return db
}

// AddArticle 新增文章到DB
func (repo *MySQLRepository) AddArticle(article *models.Article, articleContent *models.ArticleContent) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		// 新增文章
		if err := tx.Create(article).Error; err != nil {
			return err
		}

		// 设置文章内容的ArticleID
		articleContent.ArticleID = article.ID

		// 新增文章内容
		if err := tx.Create(articleContent).Error; err != nil {
			return err
		}

		return nil // 如果都添加成功，返回nil提交事务
	})
}

// ArticleExists 检查文章是否存在
func (repo *MySQLRepository) ArticleExists(articleID uint64) (bool, error) {
	var count int64
	err := repo.db.Model(&models.Article{}).Where("id = ?", articleID).Count(&count).Error
	return count > 0, err
}

// UpdateArticle 更新文章
func (repo *MySQLRepository) UpdateArticle(article *models.Article, articleContent *models.ArticleContent) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		// 更新article表中的文章
		if err := tx.Model(&models.Article{}).
			Where("id = ?", article.ID).
			Select("Title", "Picture", "Summary", "UpdatedAt").
			Updates(models.Article{
				Title:     article.Title,
				Picture:   article.Picture,
				Summary:   article.Summary,
				UpdatedAt: article.UpdatedAt,
			}).Error; err != nil {
			return err
		}

		// 更新article_content表中的文章内容
		if err := tx.Model(&models.ArticleContent{}).
			Where("article_id = ?", articleContent.ArticleID).
			Select("Content").
			Updates(models.ArticleContent{
				Content: articleContent.Content,
			}).Error; err != nil {
			return err
		}

		return nil // 如果都更新成功，返回nil提交事务
	})
}
