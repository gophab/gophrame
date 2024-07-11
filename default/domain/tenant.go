package domain

import (
	"time"

	"github.com/gophab/gophrame/core/util"

	"gorm.io/gorm"
)

type Tenant struct {
	Id               string    `gorm:"column:id;primaryKey" json:"id" primaryKey:"yes"`
	Name             string    `gorm:"column:name" json:"name"`
	Description      string    `gorm:"column:description" json:"description"`
	Logo             string    `gorm:"column:logo" json:"logo"`
	LicenseId        string    `gorm:"column:license_id" json:"licenseId"`
	Address          string    `gorm:"column:address" json:"address"`
	Telephone        string    `gorm:"column:telephone" json:"telephone"`
	Fax              string    `gorm:"column:fax" json:"fax"`
	Status           int       `gorm:"column:status" json:"status"`
	Remark           string    `gorm:"column:remark" json:"remark"`
	CreatedTime      time.Time `gorm:"column:created_time;autoCreateTime" json:"createdTime"`
	LastModifiedTime time.Time `gorm:"column:last_modified_time;autoUpdateTime" json:"lastModifiedTime"`
	CreatedBy        string    `gorm:"column:created_by" json:"createdBy"`
	LastModifiedBy   string    `gorm:"column:last_modified_by" json:"lastModifiedBy"`
	DelFlag          bool      `gorm:"column:del_flag;default:false" json:"delFlag"`
	DeletedTime      time.Time `gorm:"column:deleted_time;autoUpdateTime" json:"deleted_time"`
	DeletedBy        string    `gorm:"column:deleted_by" json:"deleted_by"`
}

func (e *Tenant) BeforeCreate(tx *gorm.DB) (err error) {
	if e.Id == "" {
		e.Id = util.UUID()
	}

	return
}

func (*Tenant) TableName() string {
	return "sys_tenant"
}
