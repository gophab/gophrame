package service

import (
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/query"
)

type EntityService[T any] interface {
	GetById(id string) (*T, error)
	GetAll(entity *T, pageable query.Pageable) (int64, []T)
	Create(entity *T) (*T, error)
	Update(entity *T) (*T, error)
	Patch(id string, column string, value any) (*T, error)
	PatchAll(id string, kv map[string]any) (*T, error)
	Delete(entity *T) error
	DeleteById(id string) error
}

type __ struct {
	RoleService       RoleService       `inject:"commonRoleService"`
	UserService       UserService       `inject:"commonUserService"`
	InviteCodeService InviteCodeService `inject:"commonInviteCodeService"`
	SocialUserService SocialUserService `inject:"commonSocialUserService"`
}

var _services = &__{}

func init() {
	inject.InjectValue("__", _services)
}
