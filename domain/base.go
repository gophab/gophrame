package domain

import (
	"time"

	"github.com/gophab/gophrame/core/snowflake"
	"github.com/gophab/gophrame/core/util"

	"gorm.io/gorm"
)

type Entity struct {
	Id               string    `gorm:"column:id;primaryKey" json:"id" primaryKey:"yes"`
	CreatedTime      time.Time `gorm:"column:created_time;autoCreateTime" json:"createdTime"`
	LastModifiedTime time.Time `gorm:"column:last_modified_time;autoUpdateTime" json:"lastModifiedTime"`
	TenantId         string    `gorm:"column:tenant_id" json:"tenantId"`
}

func (e *Entity) BeforeCreate(tx *gorm.DB) (err error) {
	if e.Id == "" {
		e.Id = util.UUID()
	}

	return
}

type AuditingEntity struct {
	Entity
	CreatedBy      string `gorm:"column:created_by" json:"createdBy"`
	LastModifiedBy string `gorm:"column:last_modified_by" json:"lastModifiedBy"`
}

type DeletableEntity struct {
	AuditingEntity
	DelFlag     bool      `gorm:"column:del_flag;default:false" json:"delFlag"`
	DeletedTime time.Time `gorm:"column:deleted_time;autoUpdateTime" json:"deleted_time"`
	DeletedBy   string    `gorm:"column:deleted_by" json:"deleted_by"`
}

type Model struct {
	Id               int64     `gorm:"primaryKey" json:"id" primaryKey:"yes"`
	CreatedTime      time.Time `gorm:"column:created_time;autoCreateTime" json:"createdTime"`
	LastModifiedTime time.Time `gorm:"column:last_modified_time;autoUpdateTime" json:"lastModifiedTime"`
	TenantId         string    `gorm:"column:tenant_id" json:"tenantId"`
}

type AuditingModel struct {
	Model
	CreatedBy      string `gorm:"created_by" json:"createdBy"`
	LastModifiedBy string `gorm:"lastModified_by" json:"lastModifiedBy"`
}

type DeletableModel struct {
	AuditingModel
	DelFlag     bool      `gorm:"column:del_flag;default:false" json:"delFlag"`
	DeletedTime time.Time `gorm:"column:deleted_time;autoUpdateTime" json:"deleted_time"`
	DeletedBy   string    `gorm:"column:deleted_by" json:"deleted_by"`
}

func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {
	if m.Id == 0 {
		m.Id = snowflake.SnowflakeIdGenerator().GetId()
	}
	return
}

type Relation struct {
	CreatedTime      time.Time `gorm:"column:created_time;autoCreateTime" json:"createdTime"`
	LastModifiedTime time.Time `gorm:"column:last_modified_time;autoUpdateTime" json:"lastModifiedTime"`
}
