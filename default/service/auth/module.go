package auth

import (
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
	"github.com/gophab/gophrame/core/util/array"
	AuthModel "github.com/gophab/gophrame/default/domain/auth"
	AuthRepository "github.com/gophab/gophrame/default/repository/auth"
	"github.com/gophab/gophrame/service"
)

type ModuleService struct {
	service.BaseService
	ModuleRepository *AuthRepository.ModuleRepository `inject:"moduleRepository"`
}

var moduleService = &ModuleService{}

func init() {
	inject.InjectValue("moduleService", moduleService)
}

func (a *ModuleService) GetById(id int64) (*AuthModel.Module, error) {
	return a.ModuleRepository.GetById(id)
}

func (a *ModuleService) GetByIds(ids []int64) ([]*AuthModel.Module, error) {
	return a.ModuleRepository.GetByIds(ids)
}

func (a *ModuleService) GetByFid(id int64) ([]*AuthModel.Module, error) {
	return a.ModuleRepository.GetByFid(id)
}

func (a *ModuleService) GetTenantModules(tenantId string) ([]*AuthModel.Module, error) {
	return a.ModuleRepository.GetTenantModules(tenantId)
}

func (a *ModuleService) SetTenantModules(tenantId string, data []*AuthModel.Module) {
	a.ModuleRepository.SetTenantModules(tenantId, data)
}

func (a *ModuleService) GetTenantModule(module, tenantId string) (*AuthModel.Module, error) {
	return a.ModuleRepository.GetTenantModule(module, tenantId)
}

func (a *ModuleService) CreateModule(module *AuthModel.Module) (*AuthModel.Module, error) {
	// 1. 获取已存在数据
	exist, err := a.GetTenantModule(module.Name, "SYSTEM")
	if err != nil {
		return nil, err
	}
	if exist != nil {
		exist.ModuleInfo = module.ModuleInfo
		return a.UpdateModule(exist)
	}

	module.Id = 0
	module.TenantId = "SYSTEM"
	_, err = a.ModuleRepository.CreateModule(module)
	if err != nil {
		return nil, err
	}

	return module, nil
}

func (a *ModuleService) UpdateModule(module *AuthModel.Module) (*AuthModel.Module, error) {
	_, err := a.ModuleRepository.UpdateModule(module)
	if err != nil {
		return nil, err
	}
	return module, nil
}

func (a *ModuleService) DeleteModule(id int64) error {
	exist, err := a.ModuleRepository.GetById(id)
	if err != nil {
		return err
	}
	if exist == nil {
		return nil
	}

	_, err = a.ModuleRepository.DeleteData(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *ModuleService) List(fid int64, title string, tenantId string, pageable query.Pageable) (int64, []*AuthModel.Module) {
	return s.ModuleRepository.List(fid, title, tenantId, pageable)
}

func (a *ModuleService) AddTenantModule(module *AuthModel.Module, tenantId string) (*AuthModel.Module, error) {
	// 1. 获取已存在数据
	exist, err := a.GetTenantModule(module.Name, tenantId)
	if err != nil {
		return nil, err
	}
	if exist != nil {
		return exist, nil
	}

	var m = &AuthModel.Module{
		ModuleInfo: module.ModuleInfo,
	}
	m.Oid = module.Id
	m.TenantId = tenantId

	if res := a.ModuleRepository.Create(m); res.Error != nil {
		return nil, res.Error
	}

	return m, nil
}

func (a *ModuleService) DeleteTenantModule(module, tenantId string) error {
	// 1. 获取已存在数据
	exist, err := a.GetTenantModule(module, tenantId)
	if err != nil {
		return err
	}
	if exist != nil {
		return nil
	}

	return a.ModuleRepository.Delete(exist).Error
}

func (s *ModuleService) MakeTree(result []*AuthModel.Module) []*AuthModel.Module {
	var menuMap = make(map[int64]*AuthModel.Module)
	for _, menu := range result {
		menu.Children = make([]*AuthModel.Module, 0)
		menuMap[menu.Id] = menu
		if menu.Oid != 0 {
			menuMap[menu.Oid] = menu
		}
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
	result, _ = array.Filter[*AuthModel.Module](result, func(m *AuthModel.Module) bool {
		return m.Fid != 0
	})
	return result
}
