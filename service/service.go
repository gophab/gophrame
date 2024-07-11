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
	Patch(id string, column string, value interface{}) (*T, error)
	PatchAll(id string, kv map[string]interface{}) (*T, error)
	Delete(entity *T) error
	DeleteById(id string) error
}

type __ struct {
	RoleService       RoleService       `inject:"roleService"`
	UserService       UserService       `inject:"userService"`
	InviteCodeService InviteCodeService `inject:"inviteCodeService"`
	SocialUserService SocialUserService `inject:"socialUserService"`
}

var _services = &__{}

func init() {
	inject.InjectValue("__", _services)
}
