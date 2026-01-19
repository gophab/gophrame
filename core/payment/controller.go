package payment

import (
	"github.com/gin-gonic/gin"

	"github.com/gophab/gophrame/core/command"
	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/eventbus"
	"github.com/gophab/gophrame/core/json"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/security"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"
	"github.com/gophab/gophrame/errors"
)

type PaymentController struct {
	controller.ResourceController
}

type PaymentRequest struct {
	Method     string `json:"method,omitempty"`
	CustomerId string `json:"customerId,omitempty"`
	BizId      string `json:"bizId,omitempty"`
	Subject    string `json:"subject,omitempty"`
	Amount     int64  `json:"amount"`
}

func (e *PaymentController) InitRouter(g *gin.RouterGroup) *gin.RouterGroup {
	api := g.Group("/openapi")
	{
		api.POST("/payment/:channel", security.HandleTokenVerify(), e.CreatePayment)
		api.PUT("/payment/:channel/:bizId/submit", security.HandleTokenVerify(), e.SubmitPayment)
		api.DELETE("/payment/:channel/:bizId", security.HandleTokenVerify(), e.CancelPayment)
		api.POST("/payment/:channel/notify", e.PaymentNotify)
	}
	return api
}

// 1. Prepay
func (e *PaymentController) CreatePayment(c *gin.Context) {
	channel, err := request.Param(c, "channel").MustString()
	if err != nil {
		logger.Error("[Payment] No payment channel: ", err.Error())
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	var paymentRequest PaymentRequest
	if err := c.BindJSON(&paymentRequest); err != nil {
		logger.Error("[Payment] Bind payment request error: ", err.Error())
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	logger.Debug("[Payment] Request to create payment: ", channel, json.String(paymentRequest))

	if command.Mode == "debug" {
		// 测试环境固定为1分钱
		paymentRequest.Amount = 1
	}

	if payment := GetPaymentChannel(channel); payment != nil {
		data, err := payment.CreatePayment(c, paymentRequest.Method, paymentRequest.Subject, paymentRequest.BizId, paymentRequest.Amount)
		if err != nil {
			logger.Error("[Payment] Create payment error: ", json.String(paymentRequest), err.Error())
			response.SystemError(c, err)
			return
		}
		response.ResponseWithJson(c, 200, data)
		return
	}

	response.NotFound(c, "Channel Not Found")
}

func (e *PaymentController) SubmitPayment(c *gin.Context) {
	channel, err := request.Param(c, "channel").MustString()
	if err != nil {
		logger.Error("[Payment] Submit payment: No payment channel - ", err.Error())
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	bizId, err := request.Param(c, "bizId").MustString()
	if err != nil {
		logger.Error("[Payment] Submit payment: No bizId - ", err.Error())
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	var paymentRequest PaymentRequest
	if err := c.BindJSON(&paymentRequest); err != nil {
		logger.Error("[Payment] Submit payment: Bind payment request error - ", err.Error())
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if command.Mode == "debug" {
		// 测试环境固定为1分钱
		paymentRequest.Amount = 1
	}

	eventbus.PublishEvent("PAYMENT_SUBMIT", &Payment{
		OrderCode: bizId,
		Channel:   channel,
		Amount:    paymentRequest.Amount,
		Status:    1, /* 提交成功 */
	})

	response.Success(c, nil)
}

func (e *PaymentController) CancelPayment(c *gin.Context) {
	channel, err := request.Param(c, "channel").MustString()
	if err != nil {
		logger.Error("[Payment] Cancel payment: No payment channel - ", err.Error())
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	bizId, err := request.Param(c, "bizId").MustString()
	if err != nil {
		logger.Error("[Payment] Cancel payment: No bizId - ", err.Error())
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if payment := GetPaymentChannel(channel); payment != nil {
		err := payment.ClosePayment(c, bizId)
		if err != nil {
			logger.Error("[Payment] Cancel payment error: ", bizId, err.Error())
			response.SystemError(c, err)
			return
		}

		eventbus.PublishEvent("PAYMENT_CANCEL", &Payment{
			OrderCode: bizId,
			Channel:   channel,
			Status:    -1, /* 取消 */
		})
		response.Success(c, nil)
		return
	}

	response.NotFound(c, "Channel Not Found")
}

func (e *PaymentController) PaymentNotify(c *gin.Context) {
	channel, err := request.Param(c, "channel").MustString()
	if err != nil {
		logger.Error("No payment channel: ", err.Error())
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	var notify = make(map[string]any)
	if err := c.ShouldBind(&notify); err != nil {
		logger.Error("Bind payment notify error: ", err.Error())
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	logger.Debug("[Payment] Received notify: ", json.String(notify))

	if payment := GetPaymentChannel(channel); payment != nil {
		data, err := payment.PaymentNotify(c, c.Request)
		if err != nil {
			logger.Error("Create payment error: ", json.String(notify), err.Error())
			if data == nil {
				// 500 错误返回
				response.SystemError(c, err)
			} else {
				// 自行构造的错误返回
				response.ResponseWithJson(c, 500, json.String(data))
			}
			return
		}

		response.Success(c, data)
		return
	}

	response.NotFound(c, "Channel Not Found")
}
