package security

import (
	"context"
	"errors"
	"strings"

	"github.com/wjshen/gophrame/core/consts"
	EmailCode "github.com/wjshen/gophrame/core/email/code"
	"github.com/wjshen/gophrame/core/logger"
	SecurityConfig "github.com/wjshen/gophrame/core/security/config"
	SecurityModel "github.com/wjshen/gophrame/core/security/model"
	SmsCode "github.com/wjshen/gophrame/core/sms/code"
	"github.com/wjshen/gophrame/core/social"
	"github.com/wjshen/gophrame/core/starter"
	"github.com/wjshen/gophrame/core/util"

	"github.com/wjshen/gophrame/security"

	"github.com/wjshen/gophrame/default/domain"
	"github.com/wjshen/gophrame/default/service"
	"github.com/wjshen/gophrame/default/service/dto"

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
	logger.Info("Initializing Default User Handler")
	security.RegisterUserHandler(new(DefaultUserHandler))
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
		return &SecurityModel.UserDetails{
			UserId:   util.StringAddr(user.Id),
			TenantId: util.StringAddr(user.TenantId),
			Login:    user.Login,
			Mobile:   user.Mobile,
			Email:    user.Email,
			Name:     user.Name,
			Avatar:   user.Avatar,
			Admin:    user.Admin,
		}, nil
	}

	return nil, errors.New("用户未注册")
}

func (h *DefaultUserHandler) GetMobileUserDetails(ctx context.Context, mobile string, code string) (*SecurityModel.UserDetails, error) {
	if h.MobileValidator == nil {
		return nil, errors.New("不支持手机验证码登录")
	}

	if !h.MobileValidator.CheckCode(h.MobileValidator, mobile, "login-pin", code) {
		return nil, errors.New("验证码不一致")
	}

	user, err := h.UserService.GetByMobile(mobile)
	if err != nil {
		return nil, err
	}

	if user == nil {
		// 是否支持自动注册
		if SecurityConfig.Setting.AutoRegister && (SecurityConfig.Setting.MobileAutoRegister == nil || *SecurityConfig.Setting.MobileAutoRegister) {
			// 用手机号/
			if user, err = h.UserService.CreateUser(&dto.User{
				Mobile:   util.StringAddr(mobile),
				Password: util.StringAddr("####*****####"),
				Status:   util.IntAddr(consts.STATUS_VALID),
			}); err != nil {
				return nil, err
			}
		}
	}

	if user != nil {
		return &SecurityModel.UserDetails{
			UserId:   util.StringAddr(user.Id),
			TenantId: util.StringAddr(user.TenantId),
			Name:     user.Name,
			Avatar:   user.Avatar,
			Mobile:   user.Mobile,
			Email:    user.Email,
			Admin:    user.Admin,
		}, nil
	}

	return nil, errors.New("该手机号用户未注册")
}

func (h *DefaultUserHandler) GetEmailUserDetails(ctx context.Context, email string, code string) (*SecurityModel.UserDetails, error) {
	if h.EmailValidator == nil {
		return nil, errors.New("不支持邮箱验证码登录")
	}

	if !h.EmailValidator.CheckCode(h.EmailValidator, email, "login-pin", code) {
		return nil, errors.New("验证码不一致")
	}

	user, err := h.UserService.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		// 是否支持自动注册
		if SecurityConfig.Setting.AutoRegister && (SecurityConfig.Setting.EmailAutoRegister == nil || *SecurityConfig.Setting.EmailAutoRegister) {
			// 用Email
			if user, err = h.UserService.CreateUser(&dto.User{
				Email:    util.StringAddr(email),
				Status:   util.IntAddr(consts.STATUS_VALID),
				Password: util.StringAddr("####*****####"),
			}); err != nil {
				return nil, err
			}
		}
	}

	if user != nil {
		return &SecurityModel.UserDetails{
			UserId:   util.StringAddr(user.Id),
			TenantId: util.StringAddr(user.TenantId),
			Name:     user.Name,
			Avatar:   user.Avatar,
			Mobile:   user.Mobile,
			Email:    user.Email,
			Admin:    user.Admin,
		}, nil
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

	result := &SecurityModel.UserDetails{
		Name:     exists.Name,
		Avatar:   exists.Avatar,
		Email:    exists.Email,
		Mobile:   exists.Mobile,
		Admin:    false,
		SocialId: util.StringAddr(exists.Id),
		TenantId: util.StringAddr(exists.TenantId),
	}

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

			if !matched && SecurityConfig.Setting.AutoRegister && (SecurityConfig.Setting.SocialAutoRegister == nil || *SecurityConfig.Setting.SocialAutoRegister) {
				// 自动注册
				user := &dto.User{
					Name:     exists.Name,
					Mobile:   exists.Mobile,
					Email:    exists.Email,
					Status:   exists.Status,
					Avatar:   exists.Avatar,
					Remark:   exists.Remark,
					Password: util.StringAddr("*****+++*****"),
				}
				if user, err := h.UserService.CreateUser(user); err == nil && user != nil {
					exists.UserId = util.StringAddr(user.Id)
				}
			}

			if exists.UserId != nil {
				h.SocialUserService.BoundSocialUser(exists.Id, *exists.UserId, exists)
			}
		} else {
			// 更新已存在的用户信息 (Merge方式）
			user := &dto.User{
				Id:     exists.UserId,
				Name:   exists.Name,
				Mobile: exists.Mobile,
				Email:  exists.Email,
				Status: exists.Status,
				Avatar: exists.Avatar,
				Remark: exists.Remark,
			}
			if u, _ := h.UserService.UpdateUser(user); u != nil {
				result.Admin = u.Admin
			}
		}
	}

	// SocialUser已生成
	if exists.UserId != nil {
		result.UserId = exists.UserId
	} else {
		result.UserId = util.StringAddr("sns:" + exists.Id)
	}
	return result, nil
}
