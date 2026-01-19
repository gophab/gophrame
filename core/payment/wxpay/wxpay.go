package wxpay

import (
	"context"
	"crypto/rsa"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gophab/gophrame/core/eventbus"
	"github.com/gophab/gophrame/core/json"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/payment"
	"github.com/gophab/gophrame/core/starter"
	"github.com/gophab/gophrame/core/util"
	"github.com/gophab/gophrame/errors"

	"github.com/gophab/gophrame/core/payment/wxpay/config"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

type WxpayService struct {
	mchPrivateKey *rsa.PrivateKey
	client        *core.Client
	handler       *notify.Handler
}

func initWxpayService() (*WxpayService, error) {
	// 加载证书
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(config.Setting.PrivateKeyFilePath)
	if err != nil {
		return nil, err
	}

	// 创建微信支付客户端
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(config.Setting.MchID, config.Setting.CertificateSerialNo, mchPrivateKey, config.Setting.APIv3Key),
	}

	client, err := core.NewClient(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(config.Setting.MchID)
	handler := notify.NewNotifyHandler(config.Setting.APIv3Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))

	return &WxpayService{
		mchPrivateKey: mchPrivateKey,
		client:        client,
		handler:       handler,
	}, nil
}

func (s *WxpayService) CreatePayment(ctx context.Context, method string, subject string, bizId string, amount int64) (string, error) {
	svc := jsapi.JsapiApiService{Client: s.client}

	openId, ok := ctx.Value("open_id").(string)
	if !ok {
		return "", errors.New(500, "No Openid")
	}

	// 构建支付请求
	req := jsapi.PrepayRequest{
		Appid:       core.String(config.Setting.AppID),
		Mchid:       core.String(config.Setting.MchID),
		Description: core.String(subject),
		OutTradeNo:  core.String(bizId),
		TimeExpire:  core.Time(time.Now()),
		NotifyUrl:   core.String(config.Setting.NotifyURL),
		Amount: &jsapi.Amount{
			Total:    core.Int64(amount), // 转换为分
			Currency: core.String("CNY"),
		},
		Payer: &jsapi.Payer{
			Openid: core.String(openId), // 需要从业务系统获取
		},
	}

	logger.Debug("[Payment] Wxpay - Prepay request: ", json.String(req))

	// 调用预支付接口
	resp, result, err := svc.Prepay(ctx, req)
	if err != nil {
		switch result.Response.StatusCode {
		case 500:
			var count = 1
			if cp := ctx.Value("__RETRY__"); cp != nil {
				count = cp.(int)
			}

			if count < 3 {
				// 间隔随机时间尝试
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
				return s.CreatePayment(context.WithValue(ctx, "__RETRY__", count+1), method, subject, bizId, amount)
			}
		case 400:

		}

		if err, b := err.(*core.APIError); b {
			switch err.Code {
			case "ORDERPAID": // 已支付
				eventbus.DispatchEvent("PAYMENT_SUCCESS",
					&payment.Payment{
						OrderCode:  bizId,
						Channel:    "wxpay",
						Amount:     amount,
						Status:     3,
						NotifyData: util.StringAddr(json.String(req)),
					})
				return "", nil
			}
		}

		logger.Error("[Payment] Wxpay - Prepay error: ", err.Error())
		return "", err
	}

	if resp == nil || resp.PrepayId == nil {
		logger.Error("[Payment] Wxpay - Prepay error: No prepay id")
		return "", errors.New(500, "Empty prepay id")
	}

	// 发送PAYMENT_CREATED event
	eventbus.PublishEvent("PAYMENT_CREATE", &payment.Payment{
		OrderCode:        bizId,
		Channel:          "wxpay",
		ChannelTradeCode: resp.PrepayId,
		Amount:           amount,
		Status:           0, /* 支付状态 -1-撤销, 0-待提交，1-提交成功，2-提交失败，3-支付成功，4-支付失败 */
		RequestData:      util.StringAddr(json.String(req)),
		ResponseData:     util.StringAddr(json.String(resp)),
	})

	// 生成前端支付所需的参数
	params := map[string]string{
		"appId":     config.Setting.AppID,
		"timeStamp": fmt.Sprintf("%d", time.Now().Unix()),
		"nonceStr":  util.GenerateRandomString(8),
		"package":   fmt.Sprintf("prepay_id=%s", util.NotNullString(resp.PrepayId)),
		"signType":  "RSA",
	}

	// 签名
	signature, err := s.signParams(params)
	if err != nil {
		logger.Error("[Payment] Wxpay - Sign error: ", err.Error())
		return "", err
	}
	params["paySign"] = signature

	// 转换为JSON返回给前端
	return json.String(params), nil
}

func (s *WxpayService) CreateRefund(ctx context.Context, paymentId string, bizId string, amount int64, reason string) (*payment.Refund, error) {
	svc := refunddomestic.RefundsApiService{Client: s.client}

	// 构建支付请求
	req := refunddomestic.CreateRequest{
		TransactionId: core.String(paymentId),
		OutTradeNo:    core.String(bizId),
		OutRefundNo:   core.String(bizId),
		Reason:        core.String(reason),
		NotifyUrl:     core.String(config.Setting.NotifyURL),
		Amount: &refunddomestic.AmountReq{
			Refund:   core.Int64(amount), // 转换为分
			Total:    core.Int64(amount), // 转换为分
			Currency: core.String("CNY"),
		},
	}

	logger.Debug("[Wxpay] Prepay request: ", json.String(req))

	// 调用预支付接口
	resp, result, err := svc.Create(ctx, req)
	if err != nil {
		if result.Response.StatusCode == 500 {
			var count = 1
			if cp := ctx.Value("__RETRY__"); cp != nil {
				count = cp.(int)
			}

			if count < 3 {
				// 间隔随机时间尝试
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
				return s.CreateRefund(context.WithValue(ctx, "__RETRY__", count+1), paymentId, bizId, amount, reason)
			}
		}
		return nil, err
	}

	if resp == nil || resp.RefundId == nil {
		return nil, errors.New(500, "Empty prepay id")
	}

	refund := &payment.Refund{
		Code:              bizId,
		OrderCode:         bizId,
		Channel:           "wxpay",
		ChannelTradeCode:  core.String(paymentId),
		ChannelRefundCode: resp.RefundId,
		Amount:            amount,
	}

	// 退款状态
	switch *resp.Status {
	case refunddomestic.STATUS_SUCCESS:
		refund.Status = 2
	case refunddomestic.STATUS_CLOSED:
		refund.Status = -1
	case refunddomestic.STATUS_PROCESSING:
		refund.Status = 1
	case refunddomestic.STATUS_ABNORMAL:
		refund.Status = -2
	}

	return refund, nil
}

func (s *WxpayService) ClosePayment(ctx context.Context, bizId string) error {
	svc := jsapi.JsapiApiService{Client: s.client}
	request := jsapi.CloseOrderRequest{
		OutTradeNo: core.String(bizId),
		Mchid:      &config.Setting.MchID,
	}
	result, err := svc.CloseOrder(ctx, request)

	if err != nil {
		logger.Error("[Wxpay] Close payment error: ", err.Error())
		return err
	}

	if result.Response.StatusCode == 200 || result.Response.StatusCode == 204 {
		// 成功
		var responseData = make([]byte, result.Response.ContentLength)
		result.Response.Body.Read(responseData)

		eventbus.DispatchEvent("PAYMENT_CLOSED", &payment.Payment{
			OrderCode:    bizId,
			RequestData:  util.StringAddr(json.String(request)),
			ResponseData: util.StringAddr(string(responseData)),
		})
		return nil
	} else {
		// 失败
		return errors.New(result.Response.StatusCode, result.Response.Status)
	}
}

func (s *WxpayService) QueryPayment(ctx context.Context, bizId string) (*payment.Payment, error) {
	payment, err := s.QueryPaymentByOrderCode(ctx, bizId)
	if err == nil && payment == nil {
		payment, err = s.QueryPaymentById(ctx, bizId)
	}

	return payment, err
}

func (s *WxpayService) QueryPaymentByOrderCode(ctx context.Context, bizId string) (*payment.Payment, error) {
	svc := jsapi.JsapiApiService{Client: s.client}
	request := jsapi.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(bizId),
		Mchid:      &config.Setting.MchID,
	}

	resp, result, err := svc.QueryOrderByOutTradeNo(ctx, request)
	if err != nil {
		if result.Response.StatusCode == 500 {
			var count = 1
			if cp := ctx.Value("__RETRY__"); cp != nil {
				count = cp.(int)
			}

			if count < 3 {
				// 间隔随机时间尝试
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
				return s.QueryPaymentByOrderCode(context.WithValue(ctx, "__RETRY__", count+1), bizId)
			}
		}
		return nil, err
	}

	if resp == nil || resp.TransactionId == nil {
		return nil, nil
	}

	p := &payment.Payment{
		OrderCode:        util.StringValue(resp.OutTradeNo),
		Channel:          "wxpay",
		ChannelTradeCode: resp.TransactionId,
		Amount:           util.Int64Value(resp.Amount.PayerTotal),
		Status:           s.paymentStatus(*resp.TradeState),
	}

	return p, nil
}

func (s *WxpayService) QueryPaymentById(ctx context.Context, bizId string) (*payment.Payment, error) {
	svc := jsapi.JsapiApiService{Client: s.client}
	request := jsapi.QueryOrderByIdRequest{
		TransactionId: core.String(bizId),
		Mchid:         &config.Setting.MchID,
	}

	resp, result, err := svc.QueryOrderById(ctx, request)
	if err != nil {
		if result.Response.StatusCode == 500 {
			var count = 1
			if cp := ctx.Value("__RETRY__"); cp != nil {
				count = cp.(int)
			}

			if count < 3 {
				// 间隔随机时间尝试
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
				return s.QueryPaymentById(context.WithValue(ctx, "__RETRY__", count+1), bizId)
			}
		}
		return nil, err
	}

	if resp == nil || resp.TransactionId == nil {
		return nil, nil
	}

	p := &payment.Payment{
		OrderCode:        util.StringValue(resp.OutTradeNo),
		Channel:          "wxpay",
		ChannelTradeCode: resp.TransactionId,
		Amount:           util.Int64Value(resp.Amount.PayerTotal),
		Status:           s.paymentStatus(*resp.TradeState),
	}

	return p, nil
}

func (s *WxpayService) PaymentNotify(ctx context.Context, request *http.Request) (any, error) {
	// Notify Sample:
	// {
	// 	"create_time": "2025-12-06T21:34:54+08:00",
	// 	"event_type": "TRANSACTION.SUCCESS",
	// 	"id": "69b2968f-283d-57e1-b53c-509920698f99",
	// 	"resource": {
	// 		"algorithm": "AEAD_AES_256_GCM",
	// 		"associated_data": "transaction",
	// 		"ciphertext": "EJp59Wt1arXWw95SnzG1GjFxivqE3OGUWCFPBqqFcxH0WlGTBWR8qJPsnK43WW51Og8XKdVYFHni7ZzyjSUom0IZvsPe+ZTX8O1MSCof63sSPtjJr19FFLxYoIJy8GISSmiAtIl90Is+bJq0l8dS3fjqdMYfHl/PWBdWOUGh2DGTwrpzTkKZN3qRs07J2jt8FJB8MwBnoLl8glOvrF6DdGqBnGTki3dCSsW3JlaybLSq6k7RhUCdbSWZqM4xwXtDGaMwlWeQgJMLehC9ZUJTUCDeNCk+Z/jQpBmOTBSyej5RkWPnJjArWSVU3H8oANNHNNsX4IFxQNuRFWaovI3yiy3gOS0996syoWgp1qQxRUFvk52EAYnoxqLQonXHPDrxcNai+brkTcnBBloaaFMNcwaw9XgpM8/50Ta+GYd1ZAdHOsmmTHgkY4DGcd0m4yimo7CcxuwT/XmjgXUmE6PUUfFy0KiURJbVXx7tthEtlg37zJpxroxcd5YRLP+UrCumkQObX3+/Ng7xRZVQBjkII2/Fdbo12nDlVdnW9hFzsM6Tmtz6u7HOOtmcSuMjsEe34LZ0Mw9Yjg==",
	// 		"nonce": "ped2seSJEuY7",
	// 		"original_type": "transaction"
	// 	},
	// 	"resource_type": "encrypt-resource",
	// 	"summary": "支付成功"
	// }
	var data map[string]any = make(map[string]any)
	req, err := s.handler.ParseNotifyRequest(ctx, request, &data)
	if err != nil {
		return map[string]any{
			"code":    "FAIL",
			"message": err.Error(),
		}, err
	}

	logger.Info("[Payment] Received notify: ", json.String(req))

	switch req.Resource.OriginalType {
	case "transaction":
		var transaction payments.Transaction
		if err := json.Json(json.String(data), &transaction); err == nil {
			var event = ""
			switch req.EventType {
			case "TRANSACTION.SUCCESS":
				event = "PAYMENT_SUCCESS"
			default:
				event = "PAYMENT_FAIL"
			}

			if event != "" {
				eventbus.DispatchEvent(event,
					&payment.Payment{
						OrderCode:        util.NotNullString(transaction.OutTradeNo),
						Channel:          "wxpay",
						ChannelTradeCode: transaction.TransactionId,
						Amount:           util.Int64Value(transaction.Amount.Total),
						Status:           s.paymentStatus(util.NotNullString(transaction.TradeState)),
						NotifyData:       util.StringAddr(json.String(req)),
					})
			}
		}
	case "refund":
		var refund refunddomestic.Refund
		if err := json.Json(json.String(data), &refund); err == nil {
			var event = ""
			switch req.EventType {
			case "REFUND.SUCCESS":
				event = "REFUND_SUCCESS"
			case "REFUND.CLOSED":
				event = "REFUND_CLOSE"
			default:
				event = "REFUND_FAIL"
			}

			eventbus.DispatchEvent(event,
				&payment.Refund{
					Code:              util.NotNullString(refund.OutTradeNo),
					OrderCode:         util.NotNullString(refund.OutTradeNo),
					Channel:           "wxpay",
					ChannelTradeCode:  refund.TransactionId,
					ChannelRefundCode: refund.RefundId,
					Amount:            util.Int64Value(refund.Amount.Refund),
					Status:            s.refundStatus(util.NotNullString((*string)(refund.Status))),
					NotifyData:        util.StringAddr(json.String(req)),
				})
		}
	}

	return nil, nil
}

func (s *WxpayService) paymentStatus(t string) int {
	switch t {
	case "SUCCESS":
		return 3
	case "REVOKED", "REFUND":
		return -1
	case "NOPAY", "USERPAYING":
		return 1
	case "CLOSED", "PAYERROR":
		return 4
	default:
		return 1
	}
}

func (s *WxpayService) refundStatus(t string) int {
	switch t {
	case "SUCCESS":
		return 2
	case "CLOSED":
		return -1
	case "PROCESSING":
		return 1
	case "ABNORMAL":
		return -2
	default:
		return 1
	}
}

func (s *WxpayService) signParams(params map[string]string) (result string, err error) {
	var content = fmt.Sprintf("%s\n%s\n%s\n%s\n", params["appId"], params["timeStamp"], params["nonceStr"], params["package"])
	var method = params["signType"]
	switch method {
	case "MD5":
		result = util.MD5(content)
	case "RSA":
		result, err = utils.SignSHA256WithRSA(content, s.mchPrivateKey)
	}
	return result, err
}

func init() {
	starter.RegisterStarter(Start)
}

func Start() {
	if config.Setting.Enabled {
		wepay, err := initWxpayService()
		if err != nil {
			logger.Error("Init Wepay Error: ", err.Error())
			return
		}
		if wepay != nil {
			payment.RegisterPaymentChannel("wxpay", wepay)
		}
	}
}
