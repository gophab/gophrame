package domain

import (
	"time"

	SecurityUtil "github.com/gophab/gophrame/core/security/util"
	"github.com/gophab/gophrame/core/snowflake"
	"github.com/gophab/gophrame/core/util"

	"gorm.io/gorm"
)

type TenantEnabled struct {
	TenantId string `gorm:"column:tenant_id" json:"tenantId"`
}

func (m *TenantEnabled) BeforeCreate(tx *gorm.DB) (err error) {
	if m.TenantId == "" {
		m.TenantId = SecurityUtil.GetCurrentTenantId(nil)
		if m.TenantId == "" {
			m.TenantId = "SYSTEM"
		}
	}
	return
}

type AuditingEnabled struct {
	CreatedBy        string     `gorm:"column:created_by;<-:create" json:"createdBy"`
	CreatedTime      time.Time  `gorm:"column:created_time;autoCreateTime;<-:create" json:"createdTime"`
	LastModifiedBy   string     `gorm:"column:last_modified_by;<-:update" json:"lastModifiedBy,omitempty"`
	LastModifiedTime *time.Time `gorm:"column:last_modified_time;autoUpdateTime;<-:update" json:"lastModifiedTime,omitempty"`
}

func (m *AuditingEnabled) BeforeCreate(tx *gorm.DB) (err error) {
	if m.CreatedBy == "" {
		m.CreatedBy = SecurityUtil.GetCurrentUserId(nil)
	}
	return
}

func (m *AuditingEnabled) BeforeSave(tx *gorm.DB) (err error) {
	if m.LastModifiedBy == "" {
		m.LastModifiedBy = SecurityUtil.GetCurrentUserId(nil)
	}
	return
}

type DeleteEnabled struct {
	DelFlag     bool       `gorm:"column:del_flag;default:false" json:"delFlag"`
	DeletedTime *time.Time `gorm:"column:deleted_time" json:"deletedTime,omitempty"`
	DeletedBy   string     `gorm:"column:deleted_by" json:"deletedBy,omitempty"`
}

func (m *DeleteEnabled) BeforeSave(tx *gorm.DB) (err error) {
	if m.DelFlag && m.DeletedBy == "" {
		m.DeletedBy = SecurityUtil.GetCurrentUserId(nil)
		if m.DeletedBy == "" {
			m.DeletedBy = "internal"
		}
		m.DeletedTime = util.TimeAddr(time.Now())
	}
	return
}

type Entity struct {
	Id string `gorm:"column:id;primaryKey" json:"id" primaryKey:"yes"`
}

func (e *Entity) BeforeCreate(tx *gorm.DB) (err error) {
	if e.Id == "" {
		e.Id = util.UUID()
	}

	return
}

type AuditingEntity struct {
	Entity
	AuditingEnabled
	TenantEnabled
}

func (e *AuditingEntity) BeforeCreate(tx *gorm.DB) (err error) {
	e.Entity.BeforeCreate(tx)
	e.AuditingEnabled.BeforeCreate(tx)
	e.TenantEnabled.BeforeCreate(tx)
	return
}

func (e *AuditingEntity) BeforeSave(tx *gorm.DB) (err error) {
	// e.Entity.BeforeSave(tx)
	e.AuditingEnabled.BeforeSave(tx)
	// e.TenantEnabled.BeforeSave(tx)
	return
}

type DeletableEntity struct {
	Entity
	AuditingEnabled
	TenantEnabled
	DeleteEnabled
}

func (e *DeletableEntity) BeforeCreate(tx *gorm.DB) (err error) {
	e.Entity.BeforeCreate(tx)
	e.AuditingEnabled.BeforeCreate(tx)
	e.TenantEnabled.BeforeCreate(tx)
	// e.DeleteEnabled.BeforeCreate(tx)
	return
}

func (e *DeletableEntity) BeforeSave(tx *gorm.DB) (err error) {
	// e.Entity.BeforeSave(tx)
	e.AuditingEnabled.BeforeSave(tx)
	// e.TenantEnabled.BeforeSave(tx)
	e.DeleteEnabled.BeforeSave(tx)
	return
}

type Model struct {
	Id int64 `gorm:"primaryKey" json:"id" primaryKey:"yes"`
}

func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {
	if m.Id == 0 {
		m.Id = snowflake.SnowflakeIdGenerator().GetId()
	}
	return
}

type AuditingModel struct {
	Model
	AuditingEnabled
	TenantEnabled
}

func (e *AuditingModel) BeforeCreate(tx *gorm.DB) (err error) {
	e.Model.BeforeCreate(tx)
	e.AuditingEnabled.BeforeCreate(tx)
	e.TenantEnabled.BeforeCreate(tx)
	return
}

func (e *AuditingModel) BeforeSave(tx *gorm.DB) (err error) {
	// e.Entity.BeforeSave(tx)
	e.AuditingEnabled.BeforeSave(tx)
	// e.TenantEnabled.BeforeSave(tx)
	return
}

type DeletableModel struct {
	Model
	AuditingEnabled
	TenantEnabled
	DeleteEnabled
}

func (e *DeletableModel) BeforeCreate(tx *gorm.DB) (err error) {
	e.Model.BeforeCreate(tx)
	e.AuditingEnabled.BeforeCreate(tx)
	e.TenantEnabled.BeforeCreate(tx)
	return
}

func (e *DeletableModel) BeforeSave(tx *gorm.DB) (err error) {
	// e.Entity.BeforeSave(tx)
	e.AuditingEnabled.BeforeSave(tx)
	e.DeleteEnabled.BeforeSave(tx)
	// e.TenantEnabled.BeforeSave(tx)
	return
}

type Relation struct {
	AuditingEnabled
}
