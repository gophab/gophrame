package service

type RoleService interface {
	GetUserRoles(userId string) []string
}

func GetRoleService() RoleService {
	return _services.RoleService
}
