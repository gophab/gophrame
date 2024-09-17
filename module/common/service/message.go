package service

import (
	"github.com/gophab/gophrame/core/cron"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
	"github.com/gophab/gophrame/core/util"

	"github.com/gophab/gophrame/module/common/domain"
	"github.com/gophab/gophrame/module/common/repository"

	"github.com/gophab/gophrame/service"
)

type MessageService struct {
	service.BaseService
	MessageRepository *repository.MessageRepository `inject:"messageRepository"`
}

var messageService = &MessageService{}

func init() {
	inject.InjectValue("messageService", messageService)

	cron.AddFunc("@every 1m", messageService.ValidateMessages)
	cron.AddFunc("@daily", messageService.HistoryMessages)
}

func (s *MessageService) GetById(id int64) (*domain.Message, error) {
	return s.MessageRepository.GetById(id)
}

func (s *MessageService) Find(conds map[string]interface{}, pageable query.Pageable) (int64, []*domain.Message, error) {
	return s.MessageRepository.Find(util.DbFields(conds), pageable)
}

func (s *MessageService) FindAvailable(conds map[string]interface{}, pageable query.Pageable) (int64, []*domain.Message, error) {
	return s.MessageRepository.FindAvailable(util.DbFields(conds), pageable)
}

func (s *MessageService) FindSimples(conds map[string]interface{}, pageable query.Pageable) (int64, []*domain.SimpleMessage, error) {
	return s.MessageRepository.FindSimples(util.DbFields(conds), pageable)
}

func (s *MessageService) FindSimplesAvailable(conds map[string]interface{}, pageable query.Pageable) (int64, []*domain.SimpleMessage, error) {
	return s.MessageRepository.FindSimplesAvailable(util.DbFields(conds), pageable)
}

func (s *MessageService) CreateMessage(message *domain.Message) (*domain.Message, error) {
	return s.MessageRepository.CreateMessage(message)
}

func (s *MessageService) PatchMessage(id int64, data map[string]interface{}) (*domain.Message, error) {
	return s.MessageRepository.PatchMessage(id, data)
}

func (s *MessageService) DeleteMessage(message *domain.Message) error {
	return s.MessageRepository.DeleteMessage(message)
}

func (s *MessageService) ValidateMessages() {
	if s.MessageRepository != nil {
		s.MessageRepository.ValidateMessages()
	}
}

func (s *MessageService) HistoryMessages() {
	if s.MessageRepository != nil {
		s.MessageRepository.HistoryMessages()
	}
}
