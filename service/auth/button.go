package auth

import (
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/query"

	AuthModel "github.com/wjshen/gophrame/domain/auth"
	"github.com/wjshen/gophrame/repository/auth"
)

type ButtonService struct {
	ButtonRepository *auth.ButtonRepository `inject:"buttonRepository"`
}

var buttonService *ButtonService = &ButtonService{}

func init() {
	inject.InjectValue("buttonService", buttonService)
}

func (s *ButtonService) GetById(id int64) (*AuthModel.Button, error) {
	return s.ButtonRepository.GetById(id)
}

func (s *ButtonService) CreateButton(data *AuthModel.Button) (bool, error) {
	return s.ButtonRepository.InsertData(data)
}

func (s *ButtonService) UpdateButton(data *AuthModel.Button) (bool, error) {
	return s.ButtonRepository.UpdateData(data)
}

func (s *ButtonService) DeleteButton(id int64) error {
	return s.ButtonRepository.DeleteData(id)
}

func (s *ButtonService) List(name string, pageable query.Pageable) (int64, []AuthModel.Button) {
	return s.ButtonRepository.List(name, pageable)
}
