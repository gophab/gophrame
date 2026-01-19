package payment

import (
	"context"
	"net/http"

	"github.com/gophab/gophrame/core"
	"github.com/gophab/gophrame/errors"
)

type Payment struct {
	Code             string  `gorm:"column:code" json:"code"`                                         /* 支付流水号 由支付系统产生 */
	OrderCode        string  `gorm:"column:order_code" json:"orderCode"`                              /* 订单号 业务系统提供 */
	Channel          string  `gorm:"column:channel" json:"channel"`                                   /* 渠道标识 */
	ChannelTradeCode *string `gorm:"column:channel_trade_code" json:"channelTradeCode"`               /* 交易流水号 由支付渠道返回 */
	Amount           int64   `gorm:"column:amount;default:0" json:"amount"`                           /* 支付金额 */
	Status           int     `gorm:"column:status;default:0" json:"status"`                           /* 支付状态 0-待提交，1-提交成功，2-提交失败，3-支付成功，4-支付失败 */
	RequestData      *string `gorm:"column:request_data;default:null" json:"requestData,omitempty"`   /* 发送数据 */
	ResponseData     *string `gorm:"column:response_data;default:null" json:"responseData,omitempty"` /* 接收数据 */
	NotifyData       *string `gorm:"column:notify_data;default:null" json:"notifyData,omitempty"`     /* 通知数据 */
}

type Refund struct {
	Code              string  `gorm:"column:code" json:"code"`                                         /* 退款流水号 由支付系统产生 */
	OrderCode         string  `gorm:"column:order_code" json:"orderCode"`                              /* 订单号 业务系统提供 */
	Channel           string  `gorm:"column:channel" json:"channel"`                                   /* 渠道标识 */
	ChannelTradeCode  *string `gorm:"column:channel_trade_code" json:"channelTradeCode"`               /* 交易流水号 由支付渠道返回 */
	ChannelRefundCode *string `gorm:"column:channel_refund_code" json:"channelRefundCode"`             /* 退款流水号 由支付渠道返回 */
	Amount            int64   `gorm:"column:amount;default:0" json:"amount"`                           /* 支付金额 */
	Status            int     `gorm:"column:status;default:0" json:"status"`                           /* 支付状态 0-待提交，1-提交成功，2-提交失败，3-支付成功，4-支付失败 */
	RequestData       *string `gorm:"column:request_data;default:null" json:"requestData,omitempty"`   /* 发送数据 */
	ResponseData      *string `gorm:"column:response_data;default:null" json:"responseData,omitempty"` /* 接收数据 */
	NotifyData        *string `gorm:"column:notify_data;default:null" json:"notifyData,omitempty"`     /* 通知数据 */
}

type PaymentChannel interface {
	CreatePayment(ctx context.Context, method string, subject string, bizId string, amount int64) (string, error)
	ClosePayment(ctx context.Context, bizId string) error
	QueryPayment(ctx context.Context, bizId string) (*Payment, error)
	CreateRefund(ctx context.Context, paymentId string, bizId string, amount int64, reason string) (*Refund, error)
	PaymentNotify(ctx context.Context, request *http.Request) (any, error)
}

var channels = make(map[string]PaymentChannel)

func RegisterPaymentChannel(code string, channel PaymentChannel) {
	channels[code] = channel
}

func GetPaymentChannel(code string) PaymentChannel {
	if result, b := channels[code]; b {
		return result
	}
	return channels[""]
}

func CreateRefund(ctx context.Context, channel string, paymentId string, bizId string, amount int64, reason string) (*Refund, error) {
	pc := GetPaymentChannel(channel)
	if pc != nil {
		return pc.CreateRefund(ctx, paymentId, bizId, amount, reason)
	}
	return nil, errors.New(404, "Channel Not Found")
}

func ClosePayment(ctx context.Context, channel string, bizId string) error {
	pc := GetPaymentChannel(channel)
	if pc != nil {
		return pc.ClosePayment(ctx, bizId)
	}
	return errors.New(404, "Channel Not Found")
}

func QueryPayment(ctx context.Context, channel string, bizId string) (*Payment, error) {
	pc := GetPaymentChannel(channel)
	if pc != nil {
		return pc.QueryPayment(ctx, bizId)
	}
	return nil, errors.New(404, "Channel Not Found")
}

func GetNotifyStatus(channel string, mode string, data string) (int, error) {
	return 0, nil
}

func GetNotifyErrorResponse(channel, mode, code string) string {
	return ""
}

func GetNotifySuccessResponse(channel, mode, code string) string {
	return ""
}

type DefaultPaymentChannel struct {
}

func (*DefaultPaymentChannel) CreatePayment(ctx context.Context, method string, subject string, bizId string, amount int64) (string, error) {
	return "", nil
}

func (*DefaultPaymentChannel) ClosePayment(ctx context.Context, paymentId string) error {
	return nil
}

func (*DefaultPaymentChannel) QueryPayment(ctx context.Context, bizId string) (*Payment, error) {
	return nil, nil
}

func (*DefaultPaymentChannel) CreateRefund(ctx context.Context, paymentId string, bizId string, amount int64, reason string) (*Refund, error) {
	return nil, nil
}

func (*DefaultPaymentChannel) PaymentNotify(ctx context.Context, request *http.Request) (any, error) {
	return core.M{}, nil
}

func (*DefaultPaymentChannel) RefundNotify(ctx context.Context, request *http.Request) (any, error) {
	return core.M{}, nil
}

var defaultPaymentChannel = &DefaultPaymentChannel{}

func init() {
	RegisterPaymentChannel("", defaultPaymentChannel)
}
