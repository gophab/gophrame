package aliyun

import (
	"errors"
	"strings"

	"github.com/gophab/gophrame/core/json"
	"github.com/gophab/gophrame/core/logger"
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
		logger.Error("[Aliyun] Send sms error: ", template, _err.Error())
		return _err
	}

	defer func() {
		if r := tea.Recover(recover()); r != nil {
			_err = r
			logger.Error("[Aliyun] Send sms error: ", template, _err.Error())
		}
	}()

	runtime := &aliyunUtil.RuntimeOptions{}
	result, _err := client.SendSmsWithOptions(s.getTemplateReq(dest, template, params), runtime)
	if _err != nil {
		logger.Error("[Aliyun] Send sms error: ", template, _err.Error())
		return _err
	}

	if *result.Body.Code != "OK" {
		_err = errors.New(result.String())
		logger.Error("[Aliyun] Send sms error: ", template, _err.Error())
	}

	return _err
}

func (s *AliyunSmsSender) getTemplateReq(phoneNumber, template string, params map[string]string) *dysmsapi20170525.SendSmsRequest {
	phoneNumber = strings.ReplaceAll(phoneNumber, " ", "")

	regionCode := "+86"
	phone := phoneNumber
	signature := config.Setting.Signature

	segs := strings.Split(phoneNumber, "-")
	if len(segs) > 1 {
		regionCode = segs[0]
		phone = segs[1]
	}

	if len(config.Setting.Templates) > 0 {
		if t, b := config.Setting.Templates[template]; b {
			template = t
		}

		if regionCode != "+86" {
			// 国际
			switch regionCode {
			case "+886", "+852", "+853": // 港澳台
				if t, b := config.Setting.Templates[template+".tc"]; b {
					template = t
				}

				if config.Setting.SignatureTC != "" {
					signature = config.Setting.SignatureTC
				}
			default: // 其他国家
				if t, b := config.Setting.Templates[template+".en"]; b {
					template = t
				}

				if config.Setting.SignatureEN != "" {
					signature = config.Setting.SignatureEN
				}
			}
		}
	}

	phoneNumber = regionCode + phone

	return &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String(signature),
		TemplateCode:  tea.String(template),
		PhoneNumbers:  tea.String(phoneNumber),
		TemplateParam: tea.String(json.String(params)),
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
