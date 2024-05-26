package qcloud

import (
	"errors"
	"strconv"

	"github.com/gophab/gophrame/core/sms/qcloud/config"
)

func CreateQcloudSmsSender() (*QcloudSmsSender, error) {
	if config.Setting.Enabled {
		return &QcloudSmsSender{}, nil
	}
	return nil, nil
}

type QcloudSmsSender struct{}

func (s *QcloudSmsSender) SendTemplateMessage(dest string, template string, params map[string]string) error {
	signature := params["signature"]
	delete(params, "signature")

	opt := NewOptions(config.Setting.AppId, config.Setting.AppKey, signature)
	opt.Debug = true

	templateCode, b := config.Setting.Templates[template]
	if b {
		if templateId, err := strconv.Atoi(templateCode); err != nil {
			return err
		} else {
			var client = NewClient(opt)
			var sm = s.getTemplateReq(dest, templateId, params)
			client.SendSMSSingle(*sm)
			return nil
		}
	}
	return errors.New("未找到对应模板")
}

func (*QcloudSmsSender) getTemplateReq(dest string, templateId int, params map[string]string) *SMSSingleReq {
	return &SMSSingleReq{
		Type:  0,
		TplID: templateId,
		Tel:   SMSTel{Nationcode: "86", Mobile: dest},
	}
}
