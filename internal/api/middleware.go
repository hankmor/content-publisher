package api

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hankmor/wechat-publisher/internal/config"
	"github.com/hankmor/wechat-publisher/internal/log"
	"go.uber.org/zap"
)

// AccessControlMiddleware 访问控制中间件
func AccessControlMiddleware(cfg *config.Config) gin.HandlerFunc {
	// 预先计算配置中所有 API Key 的哈希值，避免每次请求都进行哈希运算
	var hashedValidKeys []string
	for _, key := range cfg.API.APIKeys {
		hashedValidKeys = append(hashedValidKeys, HashAPIKey(key))
	}

	return func(ctx *gin.Context) {
		// 1. 检查 IP 白名单
		// clientIP := getClientIP(ctx)
		// if !isIPAllowed(clientIP, cfg.API.IPWhitelist) {
		// 	ctx.JSON(http.StatusForbidden, gin.H{"error": "IP 不在白名单内"})
		// 	ctx.Abort()
		// 	return
		// }

		// 2. 检查 API Key
		apiKey := strings.ToLower(getAPIKey(ctx))

		// 仅在日志中记录 API Key 的前 8 位，防止泄露
		maskedKey := "nil"
		if apiKey != "" {
			if len(apiKey) > 8 {
				maskedKey = apiKey[:8] + "..."
			} else {
				maskedKey = apiKey
			}
		}
		log.Info("API Access Control", zap.String("api_key_masked", maskedKey))

		if !isAPIKeyValid(apiKey, hashedValidKeys) {
			log.Error("invalid api key", zap.String("api_key_masked", maskedKey))
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "无效的 API Key"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// getClientIP 获取客户端 IP
func getClientIP(ctx *gin.Context) string {
	// 优先从 X-Forwarded-For 头获取（如果存在代理）
	if xForwardedFor := ctx.GetHeader("X-Forwarded-For"); xForwardedFor != "" {
		// X-Forwarded-For 格式：client, proxy1, proxy2
		parts := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(parts[0])
	}
	// 直接获取客户端 IP
	return ctx.ClientIP()
}

// isIPAllowed 检查 IP 是否在白名单内
func isIPAllowed(ip string, whitelist []string) bool {
	// 处理 IPv6 本地回环地址
	if ip == "::1" {
		ip = "127.0.0.1"
	}

	for _, allowedIP := range whitelist {
		if ip == allowedIP {
			return true
		}
	}
	return false
}

// getAPIKey 获取 API Key
func getAPIKey(ctx *gin.Context) string {
	// 从请求头获取
	if apiKey := ctx.GetHeader("X-API-Key"); apiKey != "" {
		return apiKey
	}
	// 从查询参数获取
	return ctx.Query("api_key")
}

// HashAPIKey 对 API Key 进行哈希处理（可导出）
func HashAPIKey(apiKey string) string {
	hash := sha256.Sum256([]byte(apiKey))
	return fmt.Sprintf("%x", hash)
}

// isAPIKeyValid 检查 API Key 是否有效
func isAPIKeyValid(apiKey string, hashedValidKeys []string) bool {
	if apiKey == "" {
		return false
	}

	// 客户端传递的 API Key 已经是 hash 值
	apiKeyBytes := []byte(apiKey)
	for _, validHash := range hashedValidKeys {
		// 使用常量时间比对，防止计时攻击
		if subtle.ConstantTimeCompare(apiKeyBytes, []byte(validHash)) == 1 {
			return true
		}
	}
	return false
}
