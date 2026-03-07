package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hankmor/wechat-publisher/internal/config"
)

// AccessControlMiddleware 访问控制中间件
func AccessControlMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. 检查 IP 白名单
		clientIP := getClientIP(ctx)
		if !isIPAllowed(clientIP, cfg.API.IPWhitelist) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "IP 不在白名单内"})
			ctx.Abort()
			return
		}

		// 2. 检查 API Key
		apiKey := getAPIKey(ctx)
		if !isAPIKeyValid(apiKey, cfg.API.APIKeys) {
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

// isAPIKeyValid 检查 API Key 是否有效
func isAPIKeyValid(apiKey string, validKeys []string) bool {
	for _, key := range validKeys {
		if apiKey == key {
			return true
		}
	}
	return false
}
