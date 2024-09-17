package mapper

import (
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/mapper"

	"github.com/gophab/gophrame/module/system/domain"
	"github.com/gophab/gophrame/module/system/service/dto"
)

type UserMapper struct {
}

var userMapper = &UserMapper{}

func init() {
	inject.InjectValue("userMapper", userMapper)
}

func (*UserMapper) AsDomain(src *dto.User) (dst *domain.User) {
	dst = &domain.User{}
	_ = mapper.Map(src, dst)
	return
}

func (*UserMapper) AsDomainArray(src []dto.User) (dst []domain.User) {
	if src == nil {
		return
	}

	dst = make([]domain.User, len(src))
	for index, s := range src {
		_ = mapper.Map(&s, &dst[index])
	}
	return
}

func (*UserMapper) AsDto(src *domain.User) (dst *dto.User) {
	dst = &dto.User{}
	_ = mapper.Map(src, dst)
	return
}

func (*UserMapper) AsDtoArray(src []domain.User) (dst []dto.User) {
	if src == nil {
		return
	}

	dst = make([]dto.User, len(src))
	_ = mapper.MapArrayRender(src, dst, func(s, d interface{}) {
		// s.(*domain.User)
		// d.(*dto.User)
	})

	return
}
