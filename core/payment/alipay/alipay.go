package alipay

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gophab/gophrame/core/eventbus"
	"github.com/gophab/gophrame/core/json"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/payment"
	"github.com/gophab/gophrame/core/starter"
	"github.com/gophab/gophrame/core/util"
	"github.com/gophab/gophrame/errors"

	"github.com/gophab/gophrame/core/payment/alipay/config"

	"github.com/smartwalle/alipay/v3"
)

type AlipayService struct {
	client *alipay.Client
}

func (s *AlipayService) CreatePayment(ctx context.Context, method string, subject string, bizId string, amount int64) (string, error) {
	trade := alipay.Trade{
		NotifyURL: config.Setting.NotifyURL,
		ReturnURL: config.Setting.ReturnURL,

		Subject:     subject,
		OutTradeNo:  bizId,
		TotalAmount: fmt.Sprintf("%.2f", float64(amount)/100.0),
		ProductCode: "FAST_INSTANT_TRADE_PAY",
	}
	switch method {
	case "APP":
		trade.ProductCode = "QUICK_MSECURITY_PAY"
		result, err := s.client.TradeAppPay(alipay.TradeAppPay{Trade: trade})
		if err != nil {
			return "", err
		}
		return result, nil

	case "WAP":
		trade.ProductCode = "QUICK_WAP_PAY"
		url, err := s.client.TradeWapPay(alipay.TradeWapPay{Trade: trade})
		if err != nil {
			return "", err
		}

		return url.String(), nil

	case "PAGE":
		trade.ProductCode = "FAST_INSTANT_TRADE_PAY"
		url, err := s.client.TradePagePay(alipay.TradePagePay{Trade: trade})
		if err != nil {
			return "", err
		}

		return url.String(), nil
	}

	return "", nil
}

func (s *AlipayService) CreateRefund(ctx context.Context, paymentId string, bizId string, amount int64, reason string) (*payment.Refund, error) {
	refund := alipay.TradeRefund{
		OutTradeNo:   bizId,
		OutRequestNo: bizId,
		RefundAmount: fmt.Sprintf("%.2f", float64(amount)/100.0),
		RefundReason: reason,
	}
	resp, err := s.client.TradeRefund(ctx, refund)
	if err != nil {
		return nil, err
	}

	if resp.IsSuccess() {
		return &payment.Refund{
			Code:              bizId,
			OrderCode:         bizId,
			Channel:           "alipay",
			ChannelTradeCode:  util.StringAddr(paymentId),
			ChannelRefundCode: util.StringAddr(resp.TradeNo),
			Amount:            amount,
			RequestData:       util.StringAddr(json.String(refund)),
			Status:            1,
		}, nil
	} else {
		return nil, errors.New(500, resp.Msg)
	}
}

func (s *AlipayService) ClosePayment(ctx context.Context, bizId string) error {
	close := alipay.TradeClose{
		OutTradeNo: bizId,
	}
	resp, err := s.client.TradeClose(ctx, close)
	if err != nil {
		return err
	}

	if resp.IsSuccess() {
		eventbus.DispatchEvent("PAYMENT_CLOSED", &payment.Payment{
			Code:        bizId,
			OrderCode:   bizId,
			Channel:     "alipay",
			RequestData: util.StringAddr(json.String(close)),
			Status:      -1,
		})
		return nil
	} else {
		return errors.New(500, resp.Msg)
	}
}

func (*AlipayService) QueryPayment(ctx context.Context, bizId string) (*payment.Payment, error) {
	return nil, nil
}

func (s *AlipayService) PaymentNotify(ctx context.Context, request *http.Request) (any, error) {
	return map[string]any{
		"return_code": "SUCCESS",
	}, nil
}

func initAlipayService() (*AlipayService, error) {
	client, err := alipay.New(config.Setting.AppID, config.Setting.Key, config.Setting.IsProd)
	if err != nil {
		return nil, err
	}

	return &AlipayService{
		client: client,
	}, nil
}

func init() {
	starter.RegisterStarter(Start)
}

func Start() {
	if config.Setting.Enabled {
		alipay, err := initAlipayService()
		if err != nil {
			logger.Error("Init Alipay Error: ", err.Error())
			return
		}
		if alipay != nil {
			payment.RegisterPaymentChannel("ALIPAY", alipay)
		}
	}
}
