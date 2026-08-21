package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oralhistory/oralhistory/internal/constants"
	"github.com/oralhistory/oralhistory/internal/util"
)

// Auth 校验 JWT 并注入用户上下文。
func Auth(secret string, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, fmt.Sprintf("%s: 缺少 Bearer 令牌", constants.MsgUnauthorized))
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := util.ParseToken(tokenStr, secret)
		if err != nil {
			msg := err.Error()
			if msg == "token is expired" {
				util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, "登录已过期，请重新登录")
			} else if msg == "token is malformed" {
				util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, "令牌格式错误")
			} else if msg == "token signature is invalid" {
				util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, "令牌签名无效")
			} else {
				logger.Warn("auth failed", "request_id", RequestID(c), "error", err)
				util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, fmt.Sprintf("%s: 令牌无效", constants.MsgUnauthorized))
			}
			c.Abort()
			return
		}
		c.Set(ContextKeyClaims, claims)
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyUsername, claims.Username)
		c.Set(ContextKeyRole, claims.Role)
		c.Next()
	}
}
