package domain

import (
	"time"

	"github.com/gophab/gophrame/core/util"
	"gorm.io/gorm"
)

type ShortLink struct {
	Id          string     `gorm:"column:id;primaryKey"`
	Name        string     `gorm:"column:name" json:"name"`
	Key         string     `gorm:"column:key;unique" json:"key"`
	Url         string     `gorm:"column:url" json:"url"`
	CreatedTime time.Time  `gorm:"column:created_time;autoCreateTime" json:"createdTime"`
	ExpiredTime *time.Time `gorm:"column:expired_time" json:"expiredTime"`
	TenantId    string     `gorm:"column:tenant_id;default:SYSTEM" json:"tentantId"`
	FullPath    string     `gorm:"-" json:"fullPath"`
}

func (e *ShortLink) TableName() string {
	return "sys_short_link"
}

func (e *ShortLink) BeforeCreate(tx *gorm.DB) (err error) {
	if e.Id == "" {
		e.Id = util.UUID()
	}

	return
}
