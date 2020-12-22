package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rs/xid"
)

// Entity 实体共用字段,软删除
type Entity struct {
	// ID xid 20位小写字符串全局id
	ID string `json:"id" gorm:"type:varchar(20);primary_key" form:"id"`
	// 创建时间
	CreatedAt time.Time `json:"created_at"`
	// 最后更新时间
	UpdatedAt time.Time `json:"updated_at"`
	// 软删除
	DeletedAt *time.Time `json:"-" gorm:"index"`
}

// BeforeCreate GORM hook
func (*Entity) BeforeCreate(scope *gorm.Scope) error {
	col, ok := scope.FieldByName("ID")
	if ok && col.IsBlank {
		return col.Set(xid.New().String())
	}
	return nil
}

// Addon 实体附属表共用字段，硬删除，请确认只依赖于某实体
type Addon struct {
	// ID xid 20位小写字符串全局id
	ID string `json:"id" gorm:"type:varchar(20);primary_key" form:"id"`
	// 创建时间
	CreatedAt time.Time `json:"created_at"`
	// 最后更新时间
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate GORM hook
func (*Addon) BeforeCreate(scope *gorm.Scope) error {
	col, ok := scope.FieldByName("ID")
	if ok && col.IsBlank {
		return col.Set(xid.New().String())
	}
	return nil
}

// Log 日志共用字段,不可更新和删除
type Log struct {
	// ID xid 20位小写字符串全局id
	ID string `json:"id" gorm:"type:varchar(20);primary_key" form:"id"`
	// 创建时间
	CreatedAt time.Time `json:"created_at"`
}

// BeforeCreate GORM hook
func (*Log) BeforeCreate(scope *gorm.Scope) error {
	col, ok := scope.FieldByName("ID")
	if ok && col.IsBlank {
		return col.Set(xid.New().String())
	}
	return nil
}
