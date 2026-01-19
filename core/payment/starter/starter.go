package starter

import (
	"sync"

	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/payment"
	_ "github.com/gophab/gophrame/core/payment"
	_ "github.com/gophab/gophrame/core/payment/alipay"
	_ "github.com/gophab/gophrame/core/payment/wxpay"

	"github.com/gophab/gophrame/core/payment/config"

	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/starter"
)

var (
	once sync.Once
)

func init() {
	starter.RegisterInitializor(Init)
}

func Init() {
	logger.Debug("Starting Payment: ...", config.Setting.Enabled)
	if config.Setting.Enabled {
		once.Do(func() {
			paymentController := &payment.PaymentController{}
			inject.InjectValue("paymentController", paymentController)
			controller.AddController(paymentController)
		})
	}
}
