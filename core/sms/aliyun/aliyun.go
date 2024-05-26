package aliyun

import (
	"errors"

	"github.com/gophab/gophrame/core/json"
	"github.com/gophab/gophrame/core/sms/aliyun/config"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	aliyunUtil "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
)

func CreateAliyunSmsSender() (*AliyunSmsSender, error) {
	if config.Setting.Enabled {
		return &AliyunSmsSender{}, nil
	}
	return nil, nil
}

type AliyunSmsSender struct{}

func (s *AliyunSmsSender) SendTemplateMessage(dest string, template string, params map[string]string) error {
	// 发送短信
	client, _err := s.CreateClient()
	if _err != nil {
		return _err
	}

	defer func() {
		if r := tea.Recover(recover()); r != nil {
			_err = r
		}
	}()

	runtime := &aliyunUtil.RuntimeOptions{}
	result, _err := client.SendSmsWithOptions(s.getTemplateReq(dest, template, params), runtime)
	if _err != nil {
		return _err
	}

	if *result.Body.Code != "OK" {
		_err = errors.New(result.String())
	}

	return _err
}

func (s *AliyunSmsSender) getTemplateReq(phoneNumber, template string, params map[string]string) *dysmsapi20170525.SendSmsRequest {
	signature := params["signature"]
	delete(params, "signature")

	templateCode, b := config.Setting.Templates[template]
	if b {
		return &dysmsapi20170525.SendSmsRequest{
			SignName:      tea.String(signature),
			TemplateCode:  tea.String(templateCode),
			PhoneNumbers:  tea.String(phoneNumber),
			TemplateParam: tea.String(json.String(params)),
		}
	} else {
		return nil
	}
}

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func (s *AliyunSmsSender) CreateClient() (*dysmsapi20170525.Client, error) {
	config := &openapi.Config{
		// 您的 AccessKey ID
		AccessKeyId: tea.String(config.Setting.AccessKeyId),
		// 您的 AccessKey Secret
		AccessKeySecret: tea.String(config.Setting.AccessKeySecret),
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	return dysmsapi20170525.NewClient(config)
}
