package models

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// 辅助
	DB *gorm.DB `gorm:"-"`
}

func (model *Model) AfterFind(tx *gorm.DB) (err error) {
	model.DB = tx
	return
}
func (model *Model) AfterSave(tx *gorm.DB) (err error) {
	model.DB = tx
	return
}
