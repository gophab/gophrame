package service

import (
	"fmt"
	"time"

	"github.com/gophab/gophrame/core/cron"
	"github.com/gophab/gophrame/core/eventbus"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
	"github.com/gophab/gophrame/core/util"

	"github.com/gophab/gophrame/module/common/domain"
	"github.com/gophab/gophrame/module/common/repository"

	"github.com/gophab/gophrame/service"
)

type EventService struct {
	service.BaseService
	EventRepository        *repository.EventRepository `inject:"eventRepository"`
	ContentTemplateService *ContentTemplateService     `inject:"contentTemplateService"`
}

var eventService = &EventService{}

func init() {
	inject.InjectValue("eventService", eventService)

	cron.AddFunc("@daily", eventService.HistoryEvents)

	eventbus.RegisterEventListener("SYSTEM_EVENT", eventService.TriggerEvent)
	eventbus.RegisterEventListener("ON_ACCESS_EVENT_CENTER", eventService.OnAccessEventCenter)
}

func (s *EventService) GetById(id int64) (*domain.Event, error) {
	return s.EventRepository.GetById(id)
}

func (s *EventService) Check(conds map[string]interface{}) (*domain.Event, error) {
	return s.EventRepository.Check(util.DbFields(conds))
}

func (s *EventService) Find(conds map[string]interface{}, pageable query.Pageable) (int64, []*domain.Event, error) {
	return s.EventRepository.Find(util.DbFields(conds), pageable)
}

func (s *EventService) CreateEvent(event *domain.Event) (*domain.Event, error) {
	return s.EventRepository.CreateEvent(event)
}

func (s *EventService) PatchEvent(id int64, data map[string]interface{}) (*domain.Event, error) {
	return s.EventRepository.PatchEvent(id, data)
}

func (s *EventService) DeleteEvent(event *domain.Event) error {
	return s.EventRepository.DeleteEvent(event)
}

func (s *EventService) FormatEvent(event *domain.Event) *domain.Event {
	contentTemplate, err := s.ContentTemplateService.GetByTypeAndSceneAndTenantId("event", event.Type, event.TenantId)
	if err == nil {
		return event
	}

	var params = make(map[string]interface{})
	if event.Properties != nil {
		for k, v := range *event.Properties {
			params[k] = fmt.Sprint(v)
		}
	}

	content := util.FormatParamterContentEx(contentTemplate.Content, params)
	if content != "" {
		event.Content = content
	}
	return event
}

func (s *EventService) AccessEventCenter(userId string) {
	s.EventRepository.Table("sys_event_access_log").Where("user_id=?", userId).Update("access_time", time.Now())
}

func (s *EventService) HistoryEvents() {
	if s.EventRepository != nil {
		s.EventRepository.HistoryEvents()
	}
}

func (s *EventService) TriggerEvent(event string, args ...interface{}) {
	data, b := args[0].(*domain.Event)
	if !b || data == nil {
		return
	}

	s.CreateEvent(data)
}

func (s *EventService) OnAccessEventCenter(event string, args ...interface{}) {
	userId := args[0].(string)
	if userId != "" {
		s.AccessEventCenter(userId)
	}
}
