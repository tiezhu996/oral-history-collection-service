package database

import (
	"errors"
	"fmt"

	"github.com/oralhistory/oralhistory/internal/model"
	"gorm.io/gorm"
)

// ErrMigrate 表结构迁移失败哨兵。
var ErrMigrate = errors.New("migrate failed")

// Migrate 执行表结构迁移。
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&model.User{},
		&model.Project{},
		&model.Question{},
		&model.Recording{},
		&model.TimelineMarker{},
		&model.AuditLog{},
	); err != nil {
		return fmt.Errorf("auto migrate %d tables: %w", 6, ErrMigrate)
	}
	return nil

}
