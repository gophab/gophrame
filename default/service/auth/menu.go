package auth

import (
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
	"github.com/gophab/gophrame/core/util/array"

	AuthModel "github.com/gophab/gophrame/default/domain/auth"
	AuthRepository "github.com/gophab/gophrame/default/repository/auth"

	"github.com/gophab/gophrame/service"
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
	if b, err := s.MenuRepository.CreateMenu(data); b {
		return true, nil
	} else {
		return false, err
	}
}

func (s *MenuService) UpdateMenu(data *AuthModel.Menu) (bool, error) {
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

func (s *MenuService) GetById(id int64) (*AuthModel.Menu, error) {
	result, err := s.MenuRepository.GetById(id)
	return result, err
}

func (s *MenuService) GetByFid(fid int64) ([]*AuthModel.Menu, error) {
	result, err := s.MenuRepository.GetByFid(fid)
	return result, err
}

func (s *MenuService) List(fid int64, title string, pageable query.Pageable) (int64, []*AuthModel.Menu) {
	return s.MenuRepository.List(fid, title, pageable)
}

func (s *MenuService) ListWithButtons(fid int64, title string, pageable query.Pageable) (int64, []*AuthModel.Menu) {
	return s.MenuRepository.ListWithButtons(fid, title, pageable)
}

func (s *MenuService) MakeTree(result []*AuthModel.Menu) []*AuthModel.Menu {
	var menuMap = make(map[int64]*AuthModel.Menu)
	for _, menu := range result {
		menu.Children = make([]*AuthModel.Menu, 0)
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
	result, _ = array.Filter[*AuthModel.Menu](result, func(m *AuthModel.Menu) bool {
		return m.Fid != 0
	})
	return result
}
