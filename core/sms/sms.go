package sms

type SmsSender interface {
	SendTemplateMessage(phone string, template string, params map[string]string) error
}
