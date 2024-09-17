package service

import (
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
	"github.com/gophab/gophrame/core/util/array"
	"github.com/gophab/gophrame/module/operation/domain"
	"github.com/gophab/gophrame/module/operation/repository"
	"github.com/gophab/gophrame/service"
)

type ModuleService struct {
	service.BaseService
	ModuleRepository *repository.ModuleRepository `inject:"moduleRepository"`
}

var moduleService = &ModuleService{}

func init() {
	inject.InjectValue("moduleService", moduleService)
}

func (a *ModuleService) GetById(id int64) (*domain.Module, error) {
	return a.ModuleRepository.GetById(id)
}

func (a *ModuleService) GetByIds(ids []int64) ([]*domain.Module, error) {
	return a.ModuleRepository.GetByIds(ids)
}

func (a *ModuleService) GetByFid(id int64) ([]*domain.Module, error) {
	return a.ModuleRepository.GetByFid(id)
}

func (a *ModuleService) GetTenantModules(tenantId string) ([]*domain.Module, error) {
	return a.ModuleRepository.GetTenantModules(tenantId)
}

func (a *ModuleService) SetTenantModules(tenantId string, data []*domain.Module) {
	a.ModuleRepository.SetTenantModules(tenantId, data)
}

func (a *ModuleService) GetTenantModule(module, tenantId string) (*domain.Module, error) {
	return a.ModuleRepository.GetTenantModule(module, tenantId)
}

func (a *ModuleService) CreateModule(module *domain.Module) (*domain.Module, error) {
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

func (a *ModuleService) UpdateModule(module *domain.Module) (*domain.Module, error) {
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

func (s *ModuleService) List(fid int64, title string, tenantId string, pageable query.Pageable) (int64, []*domain.Module) {
	return s.ModuleRepository.List(fid, title, tenantId, pageable)
}

func (a *ModuleService) AddTenantModule(module *domain.Module, tenantId string) (*domain.Module, error) {
	// 1. 获取已存在数据
	exist, err := a.GetTenantModule(module.Name, tenantId)
	if err != nil {
		return nil, err
	}
	if exist != nil {
		return exist, nil
	}

	var m = &domain.Module{
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

func (s *ModuleService) MakeTree(result []*domain.Module) []*domain.Module {
	var menuMap = make(map[int64]*domain.Module)
	for _, menu := range result {
		menu.Children = make([]*domain.Module, 0)
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
	result, _ = array.Filter[*domain.Module](result, func(m *domain.Module) bool {
		return m.Fid != 0
	})
	return result
}
