package mapi

import (
	"encoding/json"

	"github.com/gophab/gophrame/core/controller"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"
	"github.com/gophab/gophrame/errors"

	"github.com/gophab/gophrame/default/domain"
	"github.com/gophab/gophrame/default/service"

	"github.com/gin-gonic/gin"
)

type TenantOptionMController struct {
	controller.ResourceController
	TenantOptionService *service.SysOptionService `inject:"sysOptionService"`
}

var tenantOptionMController = &TenantOptionMController{}

func init() {
	inject.InjectValue("tenantOptionMController", tenantOptionMController)
}

func (s *TenantOptionMController) AfterInitialize() {
	s.SetResourceHandlers([]controller.ResourceHandler{
		{HttpMethod: "GET", ResourcePath: "/tenant/:id/options", Handler: s.GetTenantOptions},
		{HttpMethod: "POST", ResourcePath: "/tenant/:id/options", Handler: s.AddTenantOptions},
		{HttpMethod: "PUT", ResourcePath: "/tenant/:id/options", Handler: s.SetTenantOptions},
		{HttpMethod: "DELETE", ResourcePath: "/tenant/:id/option/:key", Handler: s.RemoveTenantOption},
		{HttpMethod: "DELETE", ResourcePath: "/tenant/:id/options", Handler: s.RemoveTenantOptions},
	})
}

func (s *TenantOptionMController) GetTenantOptions(c *gin.Context) {
	tenantId, err := request.Param(c, "id").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if tenantOptions, err := s.TenantOptionService.GetTenantOptions(tenantId); err == nil {
		result := make(map[string]string)
		for name, option := range tenantOptions.Options {
			result[name] = option.Value
		}
		response.Success(c, result)
	} else {
		response.SystemErrorMessage(c, errors.ERROR_GET_S_FAIL, err.Error())
	}
}

func (s *TenantOptionMController) AddTenantOptions(c *gin.Context) {
	tenantId, err := request.Param(c, "id").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if body, err := c.GetRawData(); err == nil {
		var data map[string]string
		_ = json.Unmarshal(body, &data)
		for k, v := range data {
			if _, err := s.TenantOptionService.AddSysOption(&domain.SysOption{
				TenantId: tenantId,
				Option: domain.Option{
					Name:      k,
					Value:     v,
					ValueType: "STRING",
				},
			}); err != nil {
				response.SystemErrorMessage(c, errors.ERROR_CREATE_FAIL, err.Error())
				return
			}
		}
		s.GetTenantOptions(c)
	} else {
		response.FailMessage(c, 400, err.Error())
	}
}

func (s *TenantOptionMController) SetTenantOptions(c *gin.Context) {
	tenantId, err := request.Param(c, "id").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if body, err := c.GetRawData(); err == nil {
		var tenantOptions = domain.SysOptions{
			TenantId: tenantId,
		}

		var data map[string]string
		_ = json.Unmarshal(body, &data)
		for k, v := range data {
			tenantOptions.Options[k] = domain.SysOption{
				TenantId: tenantId,
				Option: domain.Option{
					Name:      k,
					Value:     v,
					ValueType: "STRING",
				},
			}
		}

		if _, err := s.TenantOptionService.SetTenantOptions(&tenantOptions); err == nil {
			s.GetTenantOptions(c)
		} else {
			response.SystemErrorMessage(c, errors.ERROR_UPDATE_FAIL, err.Error())
		}
	} else {
		response.FailMessage(c, errors.INVALID_PARAMS, err.Error())
	}
}

func (s *TenantOptionMController) RemoveTenantOptions(c *gin.Context) {
	tenantId, err := request.Param(c, "id").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	if err := s.TenantOptionService.RemoveAllTenantOptions(tenantId); err != nil {
		response.SystemErrorMessage(c, errors.ERROR_DELETE_FAIL, err.Error())
	} else {
		response.Success(c, nil)
	}
}

func (s *TenantOptionMController) RemoveTenantOption(c *gin.Context) {
	tenantId, err := request.Param(c, "id").MustString()
	if err != nil {
		response.FailCode(c, errors.INVALID_PARAMS)
		return
	}

	key := c.Param("key")
	if res, err := s.TenantOptionService.RemoveTenantOption(tenantId, key); err != nil {
		response.SystemErrorMessage(c, errors.ERROR_DELETE_FAIL, err.Error())
	} else {
		response.Success(c, res)
	}
}
