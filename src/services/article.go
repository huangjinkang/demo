package services

import (
	"context"
	"demo/src/dtos"
	"demo/src/errs"
	"demo/src/models"
	"demo/src/repositories"
	"errors"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"log"
	"sync"
)

const maxSummaryLength = 200 // 文章摘要字符

type ArticleService struct {
	mysqlRepo         *repositories.MySQLRepository
	elasticsearchRepo *repositories.ElasticsearchRepository
	redisRepo         *repositories.RedisRepository
}

func NewArticleService(db *gorm.DB, esClient *elasticsearch.Client, rdb *redis.Client) *ArticleService {
	return &ArticleService{
		mysqlRepo:         repositories.NewMySQLRepository(db),
		elasticsearchRepo: repositories.NewElasticsearchRepository(esClient),
		redisRepo:         repositories.NewRedisRepository(rdb),
	}
}

// generateSummary 生成文章摘要
func generateSummary(content string) string {
	if len(content) < maxSummaryLength {
		return content
	}
	return content[:maxSummaryLength]
}

// TryLockArticle 获取文章锁
func (s *ArticleService) TryLockArticle(ctx context.Context, articleID uint64) (bool, error) {
	locked, err := s.redisRepo.LockArticleID(ctx, articleID)
	if err != nil {
		return false, err
	}
	return locked, nil
}

// UnlockArticle 释放文章锁
func (s *ArticleService) UnlockArticle(ctx context.Context, articleID uint64) error {
	// 解锁文章ID
	return s.redisRepo.UnlockArticleID(ctx, articleID)
}

// AddArticle 新增文章
func (s *ArticleService) AddArticle(ctx context.Context, articleReq *dtos.ArticleAddRequest) (articleID uint64, err error) {
	// 创建models.Article实例
	article := models.Article{
		Title:   articleReq.Title,
		Picture: articleReq.Picture,
		Summary: generateSummary(articleReq.Content),
		//CreatedAt: time.Now(),
		//UpdatedAt: time.Now(),
	}

	// 创建models.ArticleContent实例
	articleContent := models.ArticleContent{
		Content: articleReq.Content,
	}

	// 新增DB文章内容
	if err = s.mysqlRepo.AddArticle(&article, &articleContent); err != nil {
		return articleID, err
	}

	// 新增ES文章内容
	if err = s.elasticsearchRepo.AddArticle(ctx, &article); err != nil {
		/*
			TODO:
			这里需要一个更复杂的错误处理机制，如果先提交事务，这里索引增加失败时就不能回滚事务，数据就会出现仅存在数据库中；
			如果等索引添加成功时提交事务，就会出ES数据添加成功后，数据库可能没有准备好（宕机等运行异常）导致ES有数据而数据库没有，
			可以考虑使用当前的事务提交实现，在当前位置增加ES索引新增重试逻辑和错误处理队列，保证数据能正确同步到ES中。
		*/
		return articleID, err
	}

	return article.ID, nil
}

// ListArticles 获取文章列表
func (s *ArticleService) ListArticles(ctx context.Context, page, pageSize int, sortField, sortOrder string) (*dtos.ArticleListResponse, error) {
	return s.elasticsearchRepo.ListArticles(ctx, page, pageSize, sortField, sortOrder)
}

// UpdateArticle 更新文章
func (s *ArticleService) UpdateArticle(ctx context.Context, articleID uint64, articleReq *dtos.ArticleUpdateRequest) error {
	// 锁定文章
	locked, err := s.TryLockArticle(ctx, articleID)
	if err != nil {
		return err
	}
	if !locked {
		return errors.New("article update in progress, please try again later")
	}

	// 文章更新完成后解锁文章
	defer func() {
		if err := s.UnlockArticle(ctx, articleID); err != nil {
			log.Printf("Failed to unlock article with ID %d: %v", articleID, err)
		}
	}()

	// 验证文章是否存在
	exists, err := s.mysqlRepo.ArticleExists(articleID)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("article not found")
	}

	// 创建models.Article实例
	article := models.Article{
		ID:      articleID,
		Title:   articleReq.Title,
		Picture: articleReq.Picture,
		Summary: generateSummary(articleReq.Content),
		//UpdatedAt: time.Now(),
	}

	// 创建models.ArticleContent实例
	articleContent := models.ArticleContent{
		ArticleID: articleID,
		Content:   articleReq.Content,
	}

	// 使用WaitGroup等待两个异步更新操作
	var wg sync.WaitGroup
	wg.Add(2)

	// 错误通道用于从goroutines接收错误
	errChan := make(chan errs.UpdateError, 2)

	// 更新DB文章内容
	go func() {
		defer wg.Done()
		if err := s.mysqlRepo.UpdateArticle(&article, &articleContent); err != nil {
			errChan <- errs.NewUpdateError(errs.DB, err)
		}
	}()

	// 更新ES文章内容
	go func() {
		defer wg.Done()
		if err := s.elasticsearchRepo.UpdateArticle(ctx, &article); err != nil {
			errChan <- errs.NewUpdateError(errs.ES, err)
		}
	}()

	// 等待所有goroutine完成
	wg.Wait()

	// 关闭错误通道
	close(errChan)

	// 检查错误通道
	var updateErr error
	for e := range errChan {
		if e.Err != nil {
			if e.Source == errs.DB {
				// 处理数据库更新错误
			} else if e.Source == errs.ES {
				// 处理ES更新错误
			}
			log.Printf("An error occurred while updating the article! ID: %d, Source: %s, Error: %v", articleID, e.Source, e.Err)
			updateErr = e.Err
		}
	}

	if updateErr != nil {
		return updateErr
	}

	return nil
}
