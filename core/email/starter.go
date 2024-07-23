package email

import (
	"sync"

	"github.com/gophab/gophrame/core/email/config"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/starter"
)

var (
	once sync.Once
)

func init() {
	starter.RegisterStarter(Start)
}

func Start() {
	logger.Debug("Enable Email: ...", config.Setting.Enabled)
	if config.Setting.Enabled {
		once.Do(func() {
			sender := &GoEmailSender{}
			sender.Init()
			inject.InjectValue("emailSender", sender)
		})
	}
}
