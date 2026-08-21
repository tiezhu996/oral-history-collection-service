package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/oralhistory/oralhistory/internal/dto"
	"github.com/oralhistory/oralhistory/internal/model"
	"github.com/oralhistory/oralhistory/internal/service"
	"github.com/oralhistory/oralhistory/internal/util"
)

// AuditHandler 审计日志接口处理器。
type AuditHandler struct {
	auditSvc service.AuditService
	logger   *slog.Logger
}

// NewAuditHandler 构造审计处理器。
func NewAuditHandler(auditSvc service.AuditService, logger *slog.Logger) *AuditHandler {
	return &AuditHandler{auditSvc: auditSvc, logger: logger}
}

// List 审计日志列表（仅管理员）。
func (h *AuditHandler) List(c *gin.Context) {
	var p dto.PageParams
	if !bindQuery(c, &p) {
		return
	}
	p.Normalize()
	logs, total, err := h.auditSvc.List(p.Page, p.PageSize, c.Query("username"))
	if err != nil {
		c.Error(err)
		return
	}
	if logs == nil {
		logs = []model.AuditLog{}
	}
	util.OK(c, gin.H{"list": logs, "total": total, "page": p.Page, "page_size": p.PageSize})
}
// Group 按用户名聚合审计日志。
func (h *AuditHandler) Group(c *gin.Context) {
	var p dto.PageParams
	if !bindQuery(c, &p) {
		return
	}
	p.Normalize()
	logs, _, err := h.auditSvc.List(p.Page, p.PageSize, c.Query("username"))
	if err != nil {
		c.Error(err)
		return
	}
	grouped := util.GroupAuditLogs(logs)
	util.OK(c, gin.H{"groups": grouped})
}
