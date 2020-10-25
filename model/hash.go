// GORM 相关工具包
package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mitchellh/hashstructure"
	"github.com/rs/xid"
)

// Source 用作抓取数据临时存储，靠哈希值判断是否已存在记录，不可更新和删除
type Source struct {
	// ID xid 20位小写字符串全局id
	ID string `json:"id" hash:"-" gorm:"type:varchar(20);primary_key"`
	// 创建时间
	CreatedAt time.Time `json:"created_at" hash:"-"`
	// 校验和
	Hash uint64 `json:"-" hash:"-" gorm:"index"`
}

// BeforeCreate GORM hook
func (*Source) BeforeCreate(scope *gorm.Scope) error {
	col, ok := scope.FieldByName("ID")
	if ok && col.IsBlank {
		return col.Set(xid.New().String())
	}
	return nil
}

// Tracker 模型自动记录创建时间和更改时间，可用哈希快速查找业务数据是否重复。可配合一个历史表在更新和删除记录时记录历史。
type Tracker struct {
	// ID xid 20位小写字符串全局id
	ID string `json:"id" hash:"-" gorm:"type:varchar(20);primary_key"`
	// 创建时间
	CreatedAt time.Time `json:"created_at" hash:"-"`
	// 最后更新时间
	UpdatedAt time.Time `json:"updated_at" hash:"-"`
	// 校验和
	Hash uint64 `json:"-" hash:"-" gorm:"index"`
}

// BeforeCreate GORM hook
func (*Tracker) BeforeCreate(scope *gorm.Scope) error {
	col, ok := scope.FieldByName("ID")
	if ok && col.IsBlank {
		return col.Set(xid.New().String())
	}
	return nil
}

// Hash FNV hash64 of v
func Hash(v interface{}) (uint64, error) {
	return hashstructure.Hash(v, nil)
}
