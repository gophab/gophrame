package auth

import (
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/query"

	AuthModel "github.com/wjshen/gophrame/default/domain/auth"
	AuthRepository "github.com/wjshen/gophrame/default/repository/auth"

	"github.com/wjshen/gophrame/service"
)

type MenuService struct {
	service.BaseService
	MenuRepository       *AuthRepository.MenuRepository       `inject:"menuRepository"`
	MenuButtonRepository *AuthRepository.MenuButtonRepository `inject:"menuButtonRepository"`
}

var menuService *MenuService = &MenuService{}

func init() {
	inject.InjectValue("menuService", menuService)
}

func (s *MenuService) CreateMenu(data *AuthModel.Menu) (bool, error) {
	if b, err := s.MenuRepository.InsertData(data); b {
		return true, nil
	} else {
		return false, err
	}
}

func (s *MenuService) UpdateMenu(data *AuthModel.Menu) (bool, error) {
	if b, err := s.MenuRepository.UpdateData(data); b {
		return true, nil
	} else {
		return false, err
	}
}

func (s *MenuService) DeleteMenu(id int64) (bool, error) {
	if b, err := s.MenuRepository.DeleteData(id); b {
		return true, nil
	} else {
		return false, err
	}
}

func (s *MenuService) GetById(id int64) (*AuthModel.Menu, error) {
	result, err := s.MenuRepository.GetById(id)
	return &result, err
}

func (s *MenuService) GetByFid(fid int64) ([]AuthModel.Menu, error) {
	result, err := s.MenuRepository.GetByFid(fid)
	return result, err
}

func (s *MenuService) List(fid int64, title string, pageable query.Pageable) (int64, []AuthModel.Menu) {
	return s.MenuRepository.List(fid, title, pageable)
}

func (s *MenuService) ListWithButtons(fid int64, title string, pageable query.Pageable) (int64, []AuthModel.MenuWithButton) {
	return s.MenuRepository.ListWithButtons(fid, title, pageable)
}
