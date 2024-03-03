package repositories

import (
	"bytes"
	"context"
	"demo/src/dtos"
	"demo/src/models"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

const articleIndex = "article"

type articleUpdate struct {
	Title     string    `json:"title"`
	Picture   string    `json:"picture"`
	Summary   string    `json:"summary"`
	UpdatedAt time.Time `json:"updated_at"`
}

type errorResponse struct {
	Error struct {
		Type   string `json:"type"`
		Reason string `json:"reason"`
	} `json:"error"`
}

type ElasticsearchRepository struct {
	client *elasticsearch.Client
}

func NewElasticsearchRepository(client *elasticsearch.Client) *ElasticsearchRepository {
	return &ElasticsearchRepository{client: client}
}

// InitElasticsearch 初始化Elasticsearch客户端
func InitElasticsearch() *elasticsearch.Client {
	// 从.env配置文件中获取数据库连接信息
	esHost := os.Getenv("ES_HOST")

	cfg := elasticsearch.Config{
		Addresses: []string{
			esHost,
		},
	}

	// 创建ES客户端连接
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// 检查连接状态
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	fmt.Println("ES connected")
	return es
}

// AddArticle 新增ES文章
func (repo *ElasticsearchRepository) AddArticle(ctx context.Context, article *models.Article) error {
	// 将models.Article转换为JSON
	articleJSON, err := json.Marshal(article)
	if err != nil {
		return err
	}

	// 将uint64类型的ID转换为字符串
	articleID := strconv.FormatUint(article.ID, 10)

	// 创建一个Index请求
	req := esapi.IndexRequest{
		Index:      articleIndex,
		DocumentID: articleID,
		Body:       bytes.NewReader(articleJSON),
		Refresh:    "true", // 设置刷新策略
	}

	// 发送请求
	res, err := req.Do(ctx, repo.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		// 这里可以进一步解析错误信息
		return fmt.Errorf("error indexing article ID=%v", article.ID)
	}

	return nil
}

// ListArticles 获取ES文章列表
func (repo *ElasticsearchRepository) ListArticles(ctx context.Context, page, pageSize int, sortField, sortOrder string) (*dtos.ArticleListResponse, error) {
	// 设置默认排序字段和排序顺序
	if sortField == "" {
		sortField = "created_at"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}

	// 计算要跳过的文档数量
	from := (page - 1) * pageSize

	// 构建查询
	query := map[string]interface{}{
		"from": from,
		"size": pageSize,
		"sort": []interface{}{
			map[string]interface{}{
				sortField: map[string]interface{}{"order": sortOrder},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	// 执行查询
	res, err := repo.client.Search(
		repo.client.Search.WithContext(ctx),
		repo.client.Search.WithIndex("article"),
		repo.client.Search.WithBody(&buf),
		repo.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// 读取响应体内容
	responseBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading the response body: %s", err)
	}

	// 检查是否有错误
	if res.IsError() {
		var e errorResponse
		if err := json.Unmarshal(responseBytes, &e); err != nil {
			return nil, fmt.Errorf("error parsing the response body: %s", err)
		}
		return nil, fmt.Errorf("article query returns error! status: %s type: %s reason: %s", res.Status(), e.Error.Type, e.Error.Reason)
	}

	// 解析结果
	var esResponse struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source models.Article `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.Unmarshal(responseBytes, &esResponse); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %s", err)
	}

	// 计算总页数
	totalPages := int(esResponse.Hits.Total.Value) / pageSize
	if int(esResponse.Hits.Total.Value)%pageSize > 0 {
		totalPages++
	}

	// 构建文章列表
	articles := make([]models.Article, len(esResponse.Hits.Hits))
	for i, hit := range esResponse.Hits.Hits {
		articles[i] = hit.Source
	}

	// 构建响应数据
	data := dtos.ArticleListData{
		PageData: dtos.ArticleListPageData{
			Total:     esResponse.Hits.Total.Value,
			Page:      page,
			PageSize:  pageSize,
			TotalPage: totalPages,
		},
		List: articles,
	}

	response := &dtos.ArticleListResponse{
		Data:    data,
		Message: "Articles fetched successfully",
	}

	return response, nil
}

// UpdateArticle 更新ES文章
func (repo *ElasticsearchRepository) UpdateArticle(ctx context.Context, article *models.Article) error {
	// 构建ES文章更新数据
	update := struct {
		Doc articleUpdate `json:"doc"`
	}{
		Doc: articleUpdate{
			Title:     article.Title,
			Picture:   article.Picture,
			Summary:   article.Summary,
			UpdatedAt: article.UpdatedAt,
		},
	}
	articleJSON, err := json.Marshal(update)
	if err != nil {
		return err
	}

	// 将uint64类型的ID转换为字符串
	articleID := strconv.FormatUint(article.ID, 10)

	// 创建一个Update请求
	req := esapi.UpdateRequest{
		Index:      articleIndex,
		DocumentID: articleID,
		Body:       bytes.NewReader(articleJSON),
		Refresh:    "true",
	}

	// 发送请求
	res, err := req.Do(ctx, repo.client)
	fmt.Println("res:", res)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body for article ID %s: %v", articleID, err)
		}
	}(res.Body)

	if res.IsError() {
		return fmt.Errorf("error updating article! ID=%v", article.ID)
	}

	return nil
}
