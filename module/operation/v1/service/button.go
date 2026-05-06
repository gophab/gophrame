package service

import (
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"

	"github.com/gophab/gophrame/module/operation/v1/domain"
	"github.com/gophab/gophrame/module/operation/v1/repository"
)

type ButtonService struct {
	ButtonRepository *repository.ButtonRepository `inject:"buttonRepository_v1"`
}

var buttonService *ButtonService = &ButtonService{}

func init() {
	inject.InjectValue("buttonService_v1", buttonService)
}

func (s *ButtonService) GetById(id int64) (*domain.Button, error) {
	return s.ButtonRepository.GetById(id)
}

func (s *ButtonService) CreateButton(data *domain.Button) (bool, error) {
	return s.ButtonRepository.InsertData(data)
}

func (s *ButtonService) UpdateButton(data *domain.Button) (bool, error) {
	return s.ButtonRepository.UpdateData(data)
}

func (s *ButtonService) DeleteButton(id int64) error {
	return s.ButtonRepository.DeleteData(id)
}

func (s *ButtonService) List(name string, pageable query.Pageable) (int64, []domain.Button) {
	return s.ButtonRepository.List(name, pageable)
}
