package email

type EmailSender interface {
	SendTemplateEmail(addr string, template string, params map[string]string) error
}
