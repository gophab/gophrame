package captcha

import (
	"github.com/dchest/captcha"
	"github.com/gophab/gophrame/core/captcha/config"
	"github.com/gophab/gophrame/core/code"
	"github.com/mojocn/base64Captcha"
)

type Captcha struct {
	Id      string `json:"id"`
	Type    string `json:"type"`
	Image   string `json:"image"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Length  int    `json:"length"`
	Enabled bool   `json:"enabled"`
}

type CaptchaService struct {
	Store   code.CodeStore
	Driver  *base64Captcha.DriverDigit
	Captcha *base64Captcha.Captcha
}

func (s *CaptchaService) Init() {
	captcha.SetCustomStore(&CaptchaStoreAdpter{
		CodeStore: s.Store,
	})
	// 字符,公式,验证码配置
	// 生成默认数字的driver
	// cp := base64Captcha.NewCaptcha(driver, store.UseWithCtx(c))   // v8下使用redis
	s.Captcha = base64Captcha.NewCaptcha(
		base64Captcha.NewDriverDigit(config.Setting.Height, config.Setting.Width, config.Setting.Length, 0.7, 80),
		&Base64CaptchaStoreAdapter{
			CodeStore: s.Store,
		})
}

func (s *CaptchaService) Generate(gtype string) (*Captcha, error) {
	var result = &Captcha{
		Length:  config.Setting.Length,
		Width:   config.Setting.Width,
		Height:  config.Setting.Height,
		Enabled: true,
	}
	switch gtype {
	case "image":
		captchaId := captcha.NewLen(config.Setting.Length)
		result.Id = captchaId
	case "base64":
		if captchaId, b64s, _, err := s.Captcha.Generate(); err == nil {
			result.Id = captchaId
			result.Image = b64s
		}
	}
	return nil, nil
}

func (s *CaptchaService) Verify(id, value string) bool {
	return true
}
