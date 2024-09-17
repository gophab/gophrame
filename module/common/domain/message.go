package domain

import (
	"time"

	"github.com/gophab/gophrame/domain"
)

type MessageInfo struct {
	From      string     `gorm:"column:from" json:"from"`
	To        string     `gorm:"column:to" json:"to"`
	Scope     string     `gorm:"column:scope;default:TENANT" json:"scope"`
	Type      string     `gorm:"column:type;default:NOTICE" json:"type"`
	Title     string     `gorm:"column:title" json:"title" i18n:"yes"`
	Content   string     `gorm:"column:content" json:"content" i18n:"yes"`
	ValidTime *time.Time `gorm:"column:valid_time" json:"validTime"`
	DueTime   *time.Time `gorm:"column:due_time" json:"dueTime"`
	Status    int        `gorm:"column:status;default:1" json:"status"`
	Read      bool       `gorm:"column:read;->" json:"read"`
}

type Message struct {
	domain.DeletableModel
	MessageInfo
}

type SimpleMessage struct {
	domain.DeletableModel
	From      string     `gorm:"column:from" json:"from"`
	To        string     `gorm:"column:to" json:"to"`
	Scope     string     `gorm:"column:scope;default:TENANT" json:"scope"`
	Type      string     `gorm:"column:type;default:NOTICE" json:"type"`
	Title     string     `gorm:"column:title" json:"title" i18n:"yes"`
	ValidTime *time.Time `gorm:"column:valid_time" json:"validTime"`
	DueTime   *time.Time `gorm:"column:due_time" json:"dueTime"`
	Status    int        `gorm:"column:status;default:1" json:"status"`
}

func (*Message) TableName() string {
	return "sys_message"
}

type MessageHistory struct {
	Message
}

func (*MessageHistory) TableName() string {
	return "sys_message_history"
}

type MessageAccessLog struct {
	MessageId   string    `gorm:"column:message_id" json:"messageId" primaryKey:"yes"`
	UserId      string    `gorm:"column:user_id" json:"userId" primaryKey:"yes"`
	Action      string    `gorm:"column:action" json:"action" primaryKey:"yes"`
	CreatedTime time.Time `gorm:"column:created_time;autoCreateTime;<-:create" json:"createdTime"`
}

func (*MessageAccessLog) TableName() string {
	return "sys_message_access_log"
}
