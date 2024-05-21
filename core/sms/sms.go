package sms

import (
	_ "github.com/wjshen/gophrame/config"
)

type SmsSender interface {
	SendTemplateMessage(phone string, template string, params map[string]string) error
}
