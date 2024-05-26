package service

import "github.com/gophab/gophrame/core/inject"

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
