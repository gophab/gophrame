package service

import (
	"github.com/gophab/gophrame/core/eventbus"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
	SecurityUtil "github.com/gophab/gophrame/core/security/util"
	"github.com/gophab/gophrame/core/util"

	"github.com/gophab/gophrame/module/common/domain"
	"github.com/gophab/gophrame/module/common/repository"

	"github.com/gophab/gophrame/service"
)

type OperationLogService struct {
	service.BaseService
	OperationLogRepository *repository.OperationLogRepository `inject:"operationLogRepository"`
	ContentTemplateService *ContentTemplateService            `inject:"contentTemplateService"`
}

var operationLogService = &OperationLogService{}

func init() {
	inject.InjectValue("operationLogService", operationLogService)

	eventbus.RegisterEventListener("SYSTEM_LOG_OPERATION", operationLogService.logOperation)
}

func (s *OperationLogService) Find(conds map[string]any, pageable query.Pageable) (int64, []*domain.OperationLog, error) {
	count, list, err := s.OperationLogRepository.Find(util.DbFields(conds), pageable)
	if err != nil || len(list) == 0 {
		return count, list, err
	}

	// 组装
	for _, operationLog := range list {
		// i18n: template.Content
		params := make(map[string]any)
		// 1. user
		entity := GetEntity("user", operationLog.OperatorId)
		if entity != nil {
			params["operator"] = entity
		} else {
			params["operator"] = map[string]any{"id": operationLog.OperatorId}
		}

		// 2. tenant
		if operationLog.TenantId != "" && operationLog.TenantId != "SYSTEM" {
			entity = GetEntity("tenant", operationLog.TenantId)
			if entity != nil {
				params["tenant"] = entity
			} else {
				params["tenant"] = map[string]any{"id": operationLog.TenantId}
			}
		}

		// 3. target
		if operationLog.Target != "" {
			entity = GetEntity(operationLog.Target, operationLog.TargetId)
			if entity != nil {
				params["target"] = entity
			} else {
				params["target"] = map[string]any{"id": operationLog.TargetId}
			}
		}

		// 4. Location
		if operationLog.Location != "" {
			entity = GetEntity(operationLog.Location, operationLog.LocationId)
			if entity != nil {
				params["location"] = entity
			} else {
				params["location"] = map[string]any{"id": operationLog.LocationId}
			}
		}

		// 5. Content
		if operationLog.Content != "" {
			params["content"] = operationLog.Content
		}

		// 6. Time
		params["time"] = operationLog.OperatedTime.Format("2006-01-02 15:04:05")

		// Format Content
		template, err := s.ContentTemplateService.GetByTypeAndSceneAndTenantId("operation", operationLog.Operation+":"+operationLog.Target, operationLog.TenantId)
		if err == nil && template != nil {
			// 组合
			operationLog.Text = util.FormatParamterContentEx(template.Content, params)
		} else {
			operationLog.Text = util.FormatParamterContentEx(operationLog.Content, params)
		}
	}

	return count, list, err
}

func (s *OperationLogService) Append(log *domain.OperationLog) {
	if log.OperatorId == "" {
		log.OperatorId = SecurityUtil.GetCurrentUserId(nil)
	}
	if log.OperatorId == "" {
		log.OperatorId = "00000000000000000000000000000000"
	}

	if log.TenantId == "" {
		log.TenantId = SecurityUtil.GetCurrentTenantId(nil)
	}
	if log.TenantId == "" {
		log.TenantId = "SYSTEM"

		user := GetEntity("user", log.OperatorId)
		if user != nil {
			if tenantId, b := util.GetRecordField(user, "tenantId"); b {
				log.TenantId = tenantId.(string)
			}
		}
	}

	s.OperationLogRepository.Append(log)
}

func (s *OperationLogService) logOperation(event string, args ...any) {
	var operationLog = args[0].(*domain.OperationLog)
	if operationLog != nil {
		s.Append(operationLog)
	}
}
