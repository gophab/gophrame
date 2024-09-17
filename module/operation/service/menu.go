package service

import (
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
	"github.com/gophab/gophrame/core/util/array"

	"github.com/gophab/gophrame/module/operation/domain"
	"github.com/gophab/gophrame/module/operation/repository"

	"github.com/gophab/gophrame/service"
)

type MenuService struct {
	service.BaseService
	MenuRepository *repository.MenuRepository `inject:"menuRepository"`
}

var menuService *MenuService = &MenuService{}

func init() {
	inject.InjectValue("menuService", menuService)
}

func (s *MenuService) CreateMenu(data *domain.Menu) (bool, error) {
	if b, err := s.MenuRepository.CreateMenu(data); b {
		return true, nil
	} else {
		return false, err
	}
}

func (s *MenuService) UpdateMenu(data *domain.Menu) (bool, error) {
	if b, err := s.MenuRepository.UpdateMenu(data); b {
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

func (s *MenuService) GetById(id int64) (*domain.Menu, error) {
	result, err := s.MenuRepository.GetById(id)
	return result, err
}

func (s *MenuService) GetByFid(fid int64) ([]*domain.Menu, error) {
	result, err := s.MenuRepository.GetByFid(fid)
	return result, err
}

func (s *MenuService) List(fid int64, title string, pageable query.Pageable) (int64, []*domain.Menu) {
	return s.MenuRepository.List(fid, title, pageable)
}

func (s *MenuService) ListWithButtons(fid int64, title string, pageable query.Pageable) (int64, []*domain.Menu) {
	return s.MenuRepository.ListWithButtons(fid, title, pageable)
}

func (s *MenuService) MakeTree(result []*domain.Menu) []*domain.Menu {
	var menuMap = make(map[int64]*domain.Menu)
	for _, menu := range result {
		menu.Children = make([]*domain.Menu, 0)
		menuMap[menu.Id] = menu
	}
	for _, menu := range result {
		if menu.Fid != 0 {
			var parent = menuMap[menu.Fid]
			if parent != nil {
				parent.Children = append(parent.Children, menu)
			} else {
				menu.Fid = 0
			}
		}
	}
	result, _ = array.Filter[*domain.Menu](result, func(m *domain.Menu) bool {
		return m.Fid != 0
	})
	return result
}
