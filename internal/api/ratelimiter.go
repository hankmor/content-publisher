package api

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hankmor/wechat-publisher/internal/log"
	"go.uber.org/zap"
)

// RateLimitStore 频率限制存储
type RateLimitStore struct {
	mu       sync.RWMutex
	data     map[string]*APIKeyLimit // key: api_key_hash
	filePath string
}

// APIKeyLimit API Key 的限制信息
type APIKeyLimit struct {
	Hash       string    `json:"hash"`        // API Key 的哈希值
	Date       string    `json:"date"`        // 日期，格式：2006-01-02
	DraftCount int       `json:"draft_count"` // 当天创建草稿的次数
	LastReset  time.Time `json:"last_reset"`  // 上次重置时间
}

// NewRateLimitStore 创建频率限制存储
func NewRateLimitStore(dataDir string) (*RateLimitStore, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("创建数据目录失败: %w", err)
	}

	filePath := filepath.Join(dataDir, "rate_limit.json")
	store := &RateLimitStore{
		data:     make(map[string]*APIKeyLimit),
		filePath: filePath,
	}

	// 加载已有数据
	if err := store.load(); err != nil {
		log.Error("加载频率限制数据失败", zap.Error(err))
		// 如果加载失败，使用空数据继续
	}

	return store, nil
}

// load 从文件加载数据
func (s *RateLimitStore) load() error {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // 文件不存在，使用空数据
		}
		return err
	}

	var limits []APIKeyLimit
	if err := json.Unmarshal(data, &limits); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, limit := range limits {
		// 检查是否需要重置（跨天了）
		if s.needReset(&limit) {
			limit.DraftCount = 0
			limit.Date = time.Now().Format("2006-01-02")
			limit.LastReset = time.Now()
		}
		s.data[limit.Hash] = &limit
	}

	return nil
}

// save 保存数据到文件
func (s *RateLimitStore) save() error {
	s.mu.RLock()
	limits := make([]APIKeyLimit, 0, len(s.data))
	for _, limit := range s.data {
		limits = append(limits, *limit)
	}
	s.mu.RUnlock()

	data, err := json.MarshalIndent(limits, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0644)
}

// needReset 检查是否需要重置计数
func (s *RateLimitStore) needReset(limit *APIKeyLimit) bool {
	today := time.Now().Format("2006-01-02")
	return limit.Date != today
}

// CheckAndIncrement 检查并增加草稿创建次数
// 返回值：是否允许，当前次数，错误
func (s *RateLimitStore) CheckAndIncrement(apiKeyHash string, maxLimit int) (bool, int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	today := time.Now().Format("2006-01-02")
	limit, exists := s.data[apiKeyHash]

	if !exists || limit.Date != today {
		// 新的 API Key 或跨天了，创建新的记录
		limit = &APIKeyLimit{
			Hash:       apiKeyHash,
			Date:       today,
			DraftCount: 0,
			LastReset:  time.Now(),
		}
		s.data[apiKeyHash] = limit
	}

	// 检查是否超过限制
	if limit.DraftCount >= maxLimit {
		return false, limit.DraftCount, nil
	}

	// 增加计数
	limit.DraftCount++

	// 异步保存数据
	go s.save()

	return true, limit.DraftCount, nil
}

// GetCount 获取当前 API Key 的草稿创建次数
func (s *RateLimitStore) GetCount(apiKeyHash string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	today := time.Now().Format("2006-01-02")
	limit, exists := s.data[apiKeyHash]
	if !exists || limit.Date != today {
		return 0
	}
	return limit.DraftCount
}

// RateLimitConfig 频率限制配置
type RateLimitConfig struct {
	Enabled        bool // 是否启用
	MaxDraftPerDay int  // 每天最多创建草稿次数
}

// RateLimitMiddleware 频率限制中间件
func RateLimitMiddleware(store *RateLimitStore, cfg *RateLimitConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !cfg.Enabled {
			ctx.Next()
			return
		}

		// 只对创建草稿的接口进行限制
		if ctx.Request.URL.Path != "/api/draft/create" {
			ctx.Next()
			return
		}

		// 获取 API Key
		apiKey := getAPIKey(ctx)
		if apiKey == "" {
			ctx.JSON(401, gin.H{"error": "缺少 API Key"})
			ctx.Abort()
			return
		}

		// 客户端传递的 API Key 已经是 hash 值，直接使用
		apiKeyHash := apiKey

		// 检查并增加次数
		allowed, count, err := store.CheckAndIncrement(apiKeyHash, cfg.MaxDraftPerDay)
		if err != nil {
			log.Error("频率限制检查失败", zap.Error(err))
			ctx.JSON(500, gin.H{"error": "服务器内部错误"})
			ctx.Abort()
			return
		}

		if !allowed {
			ctx.JSON(429, gin.H{
				"error": fmt.Sprintf("已达到今日发布上限（%d次），请明天再试", cfg.MaxDraftPerDay),
				"current_count": count,
				"max_limit":     cfg.MaxDraftPerDay,
			})
			ctx.Abort()
			return
		}

		// 在响应头中返回剩余次数
		ctx.Header("X-RateLimit-Limit", fmt.Sprintf("%d", cfg.MaxDraftPerDay))
		ctx.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", cfg.MaxDraftPerDay-count))

		log.Info("草稿创建频率限制",
			zap.String("api_key_hash", apiKeyHash[:8]+"..."),
			zap.Int("current_count", count),
			zap.Int("max_limit", cfg.MaxDraftPerDay),
		)

		ctx.Next()
	}
}
