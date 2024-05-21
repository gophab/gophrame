package domain

import (
	"time"
)

var StandardSysOption = map[string]interface{}{
	"": "",
}

type SysOption struct {
	Option
	Public           bool      `gorm:"column:public" json:"public"`
	CreatedTime      time.Time `gorm:"column:created_time;autoCreateTime" json:"createdTime"`
	CreatedBy        string    `gorm:"column:created_by" json:"createdBy"`
	LastModifiedTime time.Time `gorm:"column:last_modified_time;autoUpdateTime" json:"lastModifiedTime"`
	LastModifiedBy   string    `gorm:"column:last_modified_by" json:"lastModifiedBy"`
	TenantId         string    `gorm:"column:tenant_id" json:"tenantId"`
}

func (*SysOption) TableName() string {
	return "sys_option"
}

type SysOptions struct {
	TenantId string
	Options  map[string]SysOption
}

func (s *SysOptions) GetOption(name string) (string, bool) {
	if option, b := s.Options[name]; b {
		return option.Value, true
	} else {
		return "", false
	}
}

func (s *SysOptions) SetOption(name string, value string) bool {
	if option, b := s.Options[name]; b {
		option.Value = value
		s.Options[name] = option
	} else {
		s.Options[name] = SysOption{
			Option: Option{
				Name:      name,
				Value:     value,
				ValueType: "STRING",
			},
			TenantId: s.TenantId,
		}
	}
	return true
}
