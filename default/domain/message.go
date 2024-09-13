package domain

import "github.com/gophab/gophrame/domain"

type MessageInfo struct {
	From string `gorm:"column:from" json:"from"`
	To   string `gorm:"column:to" json:"to"`
}

type Message struct {
	domain.AuditingModel
	MessageInfo
}
