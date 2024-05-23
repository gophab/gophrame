package service

type RoleService interface {
	LoadAllPolicy() error
}

func GetRoleService() RoleService {
	return _services.RoleService
}
