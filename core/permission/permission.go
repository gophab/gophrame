package permission

type PermissionService interface {
	CheckPermission(userId string, resourceId string, action string) (bool, error)
}
