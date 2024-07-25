package email

import (
	"crypto/tls"

	"github.com/gophab/gophrame/core/email/config"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/util"
	"gopkg.in/gomail.v2"
)

type EmailSender interface {
	SendTemplateEmail(addr string, template string, params map[string]string) error
}

type GoEmailSender struct {
	*gomail.Dialer
}

func (s *GoEmailSender) Init() {
	if config.Setting.Sender != nil {
		s.Dialer = gomail.NewDialer(
			config.Setting.Sender.Host,
			config.Setting.Sender.Port,
			config.Setting.Sender.AuthUser,
			config.Setting.Sender.AuthPass,
		)
		// 关闭SSL协议认证
		s.Dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
}

func (s *GoEmailSender) SendTemplateEmail(addr string, template string, params map[string]string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.Setting.Sender.From)
	m.SetHeader("To", addr)

	m.SetHeader("Subject", params["title"])

	m.SetBody("text/html", util.FormatParamterContent(template, params))

	if err := s.DialAndSend(m); err != nil {
		logger.Error("Send email error: ", err.Error())
		return err
	}

	return nil
}
