package service

type RoleService interface {
}

func GetRoleService() RoleService {
	return _services.RoleService
}
