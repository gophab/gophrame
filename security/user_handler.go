package security

import (
	"context"
	"errors"
	"strings"

	EmailCode "github.com/wjshen/gophrame/core/email/code"
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/logger"
	SecurityConfig "github.com/wjshen/gophrame/core/security/config"
	SecurityModel "github.com/wjshen/gophrame/core/security/model"
	SmsCode "github.com/wjshen/gophrame/core/sms/code"
	"github.com/wjshen/gophrame/core/social"
	"github.com/wjshen/gophrame/core/starter"
	"github.com/wjshen/gophrame/core/util"
	"github.com/wjshen/gophrame/service"
	"github.com/wjshen/gophrame/service/dto"

	"github.com/wjshen/gophrame/domain"

	"gorm.io/gorm"
)

type UserHandler struct {
	*gorm.DB `inject:"database"`

	MobileValidator   *SmsCode.SmsCodeValidator     `inject:"smsCodeValidator"`
	EmailValidator    *EmailCode.EmailCodeValidator `inject:"emailCodeValidator"`
	SocialUserService *service.SocialUserService    `inject:"socialUserService"`
	UserService       *service.UserService          `inject:"userService"`
}

func init() {
	starter.RegisterStarter(Start)
}

func Start() {
	logger.Info("Initializing User Handler")
	inject.InjectValue("userHandler", new(UserHandler))
}

func (h *UserHandler) GetUserDetails(ctx context.Context, username, password string) (*SecurityModel.UserDetails, error) {
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
			UserId:   user.Id,
			Password: user.Password,
			Login:    util.StringValue(user.Login),
			Mobile:   util.StringValue(user.Mobile),
			Email:    util.StringValue(user.Email),
		}, nil
	}

	return nil, errors.New("用户未注册")
}

func (h *UserHandler) GetMobileUserDetails(ctx context.Context, mobile string, code string) (*SecurityModel.UserDetails, error) {
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
				User: domain.User{
					UserBase: domain.UserBase{
						Mobile: util.StringAddr(mobile),
						Status: &domain.STATUS_VALID,
					},
					Password: "####*****####",
				},
			}); err != nil {
				return nil, err
			}
		}
	}

	if user != nil {
		return &SecurityModel.UserDetails{
			Mobile:   mobile,
			UserId:   user.Id,
			Password: user.Password,
		}, nil
	}

	return nil, errors.New("该手机号用户未注册")
}

func (h *UserHandler) GetEmailUserDetails(ctx context.Context, email string, code string) (*SecurityModel.UserDetails, error) {
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
				User: domain.User{
					UserBase: domain.UserBase{
						Email:  util.StringAddr(email),
						Status: &domain.STATUS_VALID,
					},
					Password: "####*****####",
				},
			}); err != nil {
				return nil, err
			}
		}
	}

	if user != nil {
		return &SecurityModel.UserDetails{
			Email:    email,
			UserId:   user.Id,
			Password: user.Password,
		}, nil
	}

	return nil, errors.New("该邮箱用户未注册")
}

func (h *UserHandler) getOrCreateSocialUser(socialUser domain.SocialUser) (*domain.SocialUser, error) {
	exists, err := h.SocialUserService.GetById(socialUser.Id)
	if err != nil {
		return nil, err
	}

	if exists == nil {
		// 新建
		if exists, err = h.SocialUserService.CreateSocialUser(&socialUser); err != nil {
			return nil, err
		}
	} else {
		// 更新已存在的社交账号信息
		if exists, err = h.SocialUserService.UpdateSocialUser(&socialUser); err != nil {
			return nil, err
		}
	}
	return exists, nil
}

func (h *UserHandler) GetSocialUserDetails(ctx context.Context, socialChannelId string, code string) (*SecurityModel.UserDetails, error) {
	service := social.GetSocialService(strings.Split(strings.ToLower(socialChannelId), ":")[0])
	if service == nil {
		return nil, nil
	}

	socialUser := service.GetSocialUserByCode(ctx, socialChannelId, code)
	if socialUser == nil || socialUser.Id == "" {
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
		socialUser.Id = socialUser.Type + "_" + *socialUser.OpenId
		h.getOrCreateSocialUser(*socialUser)
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
				if user, err := h.UserService.CreateUser(&dto.User{
					User: domain.User{
						UserBase: domain.UserBase{
							Name:   exists.Name,
							Mobile: exists.Mobile,
							Email:  exists.Email,
							Status: exists.Status,
							Avatar: exists.Avatar,
							Remark: exists.Remark,
						},
						Password: "*****+++*****",
					},
				}); err == nil && user != nil {
					exists.UserId = util.StringAddr(user.Id)
				}
			}

			if exists.UserId != nil {
				h.SocialUserService.BoundSocialUser(exists.Id, *exists.UserId, exists)
			}
		} else {
			// 更新已存在的用户信息 (Merge方式）
			h.UserService.UpdateUser(&dto.User{
				User: domain.User{
					UserBase: domain.UserBase{
						Entity: domain.Entity{
							Id: *exists.UserId,
						},
						Name:   exists.Name,
						Mobile: exists.Mobile,
						Email:  exists.Email,
						Status: exists.Status,
						Avatar: exists.Avatar,
						Remark: exists.Remark,
					},
				},
			})
		}
	}

	// SocialUser已生成
	if exists.UserId != nil {
		return &SecurityModel.UserDetails{
			UserId:   util.StringValue(exists.UserId),
			SocialId: exists.Id,
		}, nil
	} else {
		return &SecurityModel.UserDetails{
			UserId:   "sns:" + exists.Id,
			SocialId: exists.Id,
		}, nil
	}
}
