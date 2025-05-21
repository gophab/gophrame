package security

import (
	"context"
	"errors"
	"strings"

	EmailCode "github.com/gophab/gophrame/core/email/code"
	"github.com/gophab/gophrame/core/logger"
	SecurityModel "github.com/gophab/gophrame/core/security/model"
	SmsCode "github.com/gophab/gophrame/core/sms/code"
	"github.com/gophab/gophrame/core/social"
	"github.com/gophab/gophrame/core/starter"
	"github.com/gophab/gophrame/core/util"
	"github.com/gophab/gophrame/core/util/array"

	"github.com/gophab/gophrame/security"

	"github.com/gophab/gophrame/module/system/domain"
	"github.com/gophab/gophrame/module/system/service"

	"gorm.io/gorm"
)

type DefaultUserHandler struct {
	*gorm.DB          `inject:"database"`
	MobileValidator   *SmsCode.SmsCodeValidator     `inject:"smsCodeValidator"`
	EmailValidator    *EmailCode.EmailCodeValidator `inject:"emailCodeValidator"`
	SocialUserService *service.SocialUserService    `inject:"socialUserService"`
	UserService       *service.UserService          `inject:"userService"`
	security.UserHandler
}

func init() {
	starter.RegisterStarter(Start)
}

func Start() {
	logger.Debug("Initializing Default User Handler")
	security.RegisterUserHandler(new(DefaultUserHandler))
}

func User2UserDetails(user *domain.User) *SecurityModel.UserDetails {
	if user != nil {
		return &SecurityModel.UserDetails{
			UserId:   util.StringAddr(user.Id),
			TenantId: util.StringAddr(user.TenantId),
			Login:    user.Login,
			Mobile:   user.Mobile,
			Email:    user.Email,
			Name:     user.Name,
			Avatar:   user.Avatar,
			Admin:    user.Admin,
			Roles: array.Map(user.Roles, func(item *domain.Role) string {
				return item.Name
			}),
		}
	} else {
		return nil
	}
}

func SocialUser2UserDetails(exists *domain.SocialUser) *SecurityModel.UserDetails {
	if exists != nil {
		return &SecurityModel.UserDetails{
			Name:     exists.Name,
			Avatar:   exists.Avatar,
			Email:    exists.Email,
			Mobile:   exists.Mobile,
			Admin:    false,
			SocialId: util.StringAddr(exists.Id),
			TenantId: util.StringAddr(exists.TenantId),
		}
	} else {
		return nil
	}
}

func (h *DefaultUserHandler) GetUserDetails(ctx context.Context, username, password string) (*SecurityModel.UserDetails, error) {
	user, err := h.UserService.Get(username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("用户未注册")
	}

	if user.Password != util.SHA1(password) {
		return nil, errors.New("用户名或密码错误")
	}

	if user.Id != "" {
		return User2UserDetails(user), nil
	}

	return nil, errors.New("用户未注册")
}

func (h *DefaultUserHandler) GetMobileUserDetails(ctx context.Context, mobile string, code string) (*SecurityModel.UserDetails, error) {
	if h.MobileValidator == nil {
		return nil, errors.New("不支持手机验证码登录")
	}

	if !h.MobileValidator.CheckCode(h.MobileValidator, mobile, "login", code) {
		return nil, errors.New("验证码不一致")
	}

	user, err := h.UserService.GetByMobile(mobile)
	if user == nil {
		mobile = strings.Replace(mobile, "+86-", "", -1)
		user, err = h.UserService.GetByMobile(mobile)
	}

	if err != nil {
		return nil, err
	}

	if user != nil && user.Id != "" {
		return User2UserDetails(user), nil
	}

	return nil, errors.New("该手机号用户未注册")
}

func (h *DefaultUserHandler) GetEmailUserDetails(ctx context.Context, email string, code string) (*SecurityModel.UserDetails, error) {
	if h.EmailValidator == nil {
		return nil, errors.New("不支持邮箱验证码登录")
	}

	if !h.EmailValidator.CheckCode(h.EmailValidator, email, "login", code) {
		return nil, errors.New("验证码不一致")
	}

	user, err := h.UserService.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	if user != nil && user.Id != "" {
		return User2UserDetails(user), nil
	}

	return nil, errors.New("该邮箱用户未注册")
}

func (h *DefaultUserHandler) getOrCreateSocialUser(socialUser social.SocialUser) (*domain.SocialUser, error) {
	exists, err := h.SocialUserService.GetById(socialUser.GetId())
	if err != nil {
		return nil, err
	}

	if exists == nil {
		// 新建
		exists = &domain.SocialUser{
			SocialUser: socialUser,
		}
		if exists, err = h.SocialUserService.CreateSocialUser(exists); err != nil {
			return nil, err
		}
	} else {
		// 更新已存在的社交账号信息
		exists.SocialUser = socialUser
		if exists, err = h.SocialUserService.UpdateSocialUser(exists); err != nil {
			return nil, err
		}
	}
	return exists, nil
}

func (h *DefaultUserHandler) GetSocialUserDetails(ctx context.Context, socialChannelId string, code string) (*SecurityModel.UserDetails, error) {
	service := social.GetSocialService(strings.Split(strings.ToLower(socialChannelId), ":")[0])
	if service == nil {
		return nil, nil
	}

	socialUser := service.GetSocialUserByCode(ctx, socialChannelId, code)
	if socialUser == nil || socialUser.UserId == nil {
		return nil, nil
	}

	// 1. OpenId(+) SocialId(-) UserId(-)
	// 2. OpenId(+) SocialId(+) UserId(-)
	// 3. OpenId(+) SocialId(+) UserId(+)
	exists, err := h.getOrCreateSocialUser(*socialUser)
	if err != nil {
		return nil, err
	}

	if socialUser.SocialId != nil && socialUser.OpenId != nil {
		// create sub
		h.getOrCreateSocialUser(*socialUser)
	}

	result := SocialUser2UserDetails(exists)

	if exists.SocialId != nil {
		// 确定社交平台唯一账户
		if exists.UserId == nil {
			var matched = false
			if !matched && exists.Mobile != nil {
				if user, err := h.UserService.GetByMobile(*exists.Mobile); err == nil && user != nil {
					exists.UserId = util.StringAddr(user.Id)
					matched = true
				}
			}

			if !matched && exists.Email != nil {
				if user, err := h.UserService.GetByEmail(*exists.Email); err == nil && user != nil {
					exists.UserId = util.StringAddr(user.Id)
					matched = true
				}
			}

			if exists.UserId != nil {
				h.SocialUserService.BoundSocialUser(exists.Id, *exists.UserId, exists)
			}
		}
	}

	// SocialUser已生成
	if exists.UserId != nil {
		if user, err := h.UserService.GetById(*exists.UserId); err == nil {
			// use User information
			result = User2UserDetails(user)
		}
		result.UserId = exists.UserId
	} else {
		result.UserId = util.StringAddr("sns:" + exists.Id)
	}
	return result, nil
}
