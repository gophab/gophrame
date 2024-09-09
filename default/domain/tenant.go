package domain

import (
	"time"

	SecurityUtil "github.com/gophab/gophrame/core/security/util"
	"github.com/gophab/gophrame/core/util"

	"gorm.io/gorm"
)

type Tenant struct {
	Id               string    `gorm:"column:id;primaryKey" json:"id" primaryKey:"yes"`
	Name             string    `gorm:"column:name" json:"name" i18n:"yes"`
	NameCn           string    `gorm:"column:name_cn;default:null" json:"nameCn,omitempty"`
	NameTw           string    `gorm:"column:name_tw;default:null" json:"nameTw,omitempty"`
	NameEn           string    `gorm:"column:name_en;default:null" json:"nameEn,omitempty"`
	Description      string    `gorm:"column:description;default:null" json:"description,omitempty"`
	Logo             string    `gorm:"column:logo;default:null" json:"logo,omitempty"`
	LicenseId        string    `gorm:"column:license_id;default:null" json:"licenseId,omitempty"`
	Address          string    `gorm:"column:address;default:null" json:"address,omitempty"`
	Telephone        string    `gorm:"column:telephone;default:null" json:"telephone,omitempty"`
	Fax              string    `gorm:"column:fax;default:null" json:"fax,omitempty"`
	Status           int       `gorm:"column:status;default:0" json:"status"`
	Remark           string    `gorm:"column:remark;default:null" json:"remark,omitempty"`
	CreatedTime      time.Time `gorm:"column:created_time;autoCreateTime" json:"createdTime"`
	LastModifiedTime time.Time `gorm:"column:last_modified_time;autoUpdateTime" json:"lastModifiedTime"`
	CreatedBy        string    `gorm:"column:created_by" json:"createdBy,omitempty"`
	LastModifiedBy   string    `gorm:"column:last_modified_by" json:"lastModifiedBy,omitempty"`
	DelFlag          bool      `gorm:"column:del_flag;default:false" json:"delFlag"`
	DeletedTime      time.Time `gorm:"column:deleted_time;autoUpdateTime" json:"deleted_time,omitempty"`
	DeletedBy        string    `gorm:"column:deleted_by;default:null" json:"deleted_by,omitempty"`
}

func (e *Tenant) BeforeCreate(tx *gorm.DB) (err error) {
	if e.Id == "" {
		e.Id = util.UUID()
	}
	if e.CreatedBy == "" {
		e.CreatedBy = SecurityUtil.GetCurrentUserId(nil)
	}
	return
}

func (e *Tenant) BeforeSave(tx *gorm.DB) (err error) {
	if e.LastModifiedBy == "" {
		e.LastModifiedBy = SecurityUtil.GetCurrentUserId(nil)
	}
	return
}

func (*Tenant) TableName() string {
	return "sys_tenant"
}
