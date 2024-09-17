package domain

import (
	"fmt"
	"reflect"
	"time"

	SecurityUtil "github.com/gophab/gophrame/core/security/util"
	"github.com/gophab/gophrame/domain"
)

type OperationLog struct {
	domain.Model
	OperatorId   string             `gorm:"column:operator_id" json:"operatorId"`
	Operation    string             `gorm:"column:operation" json:"operation"`
	Target       string             `gorm:"column:target" json:"target"`
	TargetId     string             `gorm:"column:target_id" json:"targetId"`
	Location     string             `gorm:"column:location" json:"location"`
	LocationId   string             `gorm:"column:location_id" json:"locationId"`
	Content      string             `gorm:"column:content" json:"content"`
	OperatedTime time.Time          `gorm:"column:operatedTime;autoCreateTime;->" json:"operatedTime"`
	Properties   *domain.Properties `gorm:"column:properties;type:json" json:"properties,omitempty"`
	TenantId     string             `gorm:"column:tenant_id" json:"tenantId"`
	Text         string             `gorm:"-" json:"text"` /* 组合文本 */
}

func (o *OperationLog) WithTarget(target string, targetId interface{}) *OperationLog {
	o.Target = target
	o.TargetId = fmt.Sprint(targetId)
	return o
}

func (o *OperationLog) WithTargetEx(target interface{}) *OperationLog {
	if target != nil {
		v := reflect.Indirect(reflect.ValueOf(target))
		o.Target = v.Type().Name()
		if id := v.FieldByName("Id"); !id.IsZero() {
			o.TargetId = fmt.Sprint(id)
		}
	}

	return o
}

func (o *OperationLog) WithLocation(location string, locationId interface{}) *OperationLog {
	o.Location = location
	o.LocationId = fmt.Sprint(locationId)
	return o
}

func (o *OperationLog) WithLocationEx(location interface{}) *OperationLog {
	if location != nil {
		v := reflect.Indirect(reflect.ValueOf(location))
		o.Location = v.Type().Name()
		if id := v.FieldByName("Id"); !id.IsZero() {
			o.LocationId = fmt.Sprint(id)
		}
	}
	return o
}

func (o *OperationLog) WithContent(content string) *OperationLog {
	o.Content = content
	return o
}

func (OperationLog) TableName() string {
	return "sys_operation_log"
}

func NewOperationLog(operation string) *OperationLog {
	return &OperationLog{
		OperatorId:   SecurityUtil.GetCurrentUserId(nil),
		TenantId:     SecurityUtil.GetCurrentTenantId(nil),
		Operation:    operation,
		OperatedTime: time.Now(),
	}
}
