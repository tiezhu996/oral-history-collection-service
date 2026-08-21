package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
			switch {
			case errors.Is(err, jwt.ErrTokenExpired):
				logger.Warn("token expired", "request_id", RequestID(c), "error", err)
				util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, "登录已过期，请重新登录")
			case errors.Is(err, jwt.ErrTokenMalformed):
				logger.Warn("token malformed", "request_id", RequestID(c), "error", err)
				util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, "令牌格式错误")
			case errors.Is(err, jwt.ErrTokenSignatureInvalid):
				logger.Warn("token signature invalid", "request_id", RequestID(c), "error", err)
				util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, "令牌签名无效")
			default:
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
