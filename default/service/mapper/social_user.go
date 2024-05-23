package mapper

import (
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/mapper"
	"github.com/wjshen/gophrame/core/mapper/converter"

	"github.com/wjshen/gophrame/default/domain"
	"github.com/wjshen/gophrame/default/service/dto"
)

type SocialUserMapper struct {
}

var socialUserMapperOption = &converter.StructOption{
	BannedFields: []string{"Id"},
}

var socialUserMapper = &SocialUserMapper{}

func init() {
	inject.InjectValue("socialUserMapper", socialUserMapper)
}

func (*SocialUserMapper) AsDomain(src *dto.User) (dst *domain.SocialUser) {
	dst = &domain.SocialUser{}
	_ = mapper.MapOption(src, dst, socialUserMapperOption)
	return
}

func (*SocialUserMapper) AsDomainArray(src []dto.User) (dst []domain.SocialUser) {
	if src == nil {
		return
	}

	dst = make([]domain.SocialUser, len(src))
	for index, s := range src {
		_ = mapper.MapOption(&s, &dst[index], socialUserMapperOption)
	}
	return
}

func (*SocialUserMapper) AsDto(src *domain.SocialUser) (dst *dto.User) {
	dst = &dto.User{}
	_ = mapper.MapOptionRender(src, dst, socialUserMapperOption, func(s, d interface{}) {
		// s.(*domain.User)
		// d.(*dto.User)
	})
	return
}

func (*SocialUserMapper) AsDtoArray(src []domain.SocialUser) (dst []dto.User) {
	if src == nil {
		return
	}

	dst = make([]dto.User, len(src))
	_ = mapper.MapArrayOptionRender(src, dst, socialUserMapperOption, func(s, d interface{}) {
		// s.(*domain.User)
		// d.(*dto.User)
	})

	return
}
