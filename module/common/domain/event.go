package domain

import (
	"time"

	"github.com/gophab/gophrame/domain"
)

type EventInfo struct {
	Source     string             `gorm:"column:source" json:"source"`
	SourceId   string             `gorm:"column:source_id" json:"sourceId"`
	Type       string             `gorm:"column:type" json:"type"`
	Target     string             `gorm:"column:target" json:"target"`
	Scope      string             `gorm:"column:scope" json:"scope"`
	Content    string             `gorm:"column:content" json:"content"`
	Status     int                `gorm:"column:status;default:1" json:"status"`
	Properties *domain.Properties `gorm:"column:properties;type:json" json:"properties,omitempty"`
}

type Event struct {
	domain.DeletableModel
	EventInfo
	Read bool `gorm:"-" json:"read"`
}

func (*Event) TableName() string {
	return "sys_event"
}

type EventHistory struct {
	Event
}

func (*EventHistory) TableName() string {
	return "sys_event_history"
}

type EventAccessLog struct {
	UserId     string    `gorm:"column:user_id;primaryKey" json:"userId" primaryKey:"yes"`
	Action     string    `gorm:"column:action;primaryKey" json:"action" primaryKey:"yes"`
	AccessTime time.Time `gorm:"column:access_time;autoCreateTime;autoUpdateTime" json:"accessTime"`
}

func (*EventAccessLog) TableName() string {
	return "sys_event_access_log"
}
