package service

import (
	"strings"
	"time"

	"github.com/gophab/gophrame/core"
	"github.com/gophab/gophrame/core/consts"
	"github.com/gophab/gophrame/core/eventbus"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/security"
	SecurityUtil "github.com/gophab/gophrame/core/security/util"
	"github.com/gophab/gophrame/core/util"
	"github.com/gophab/gophrame/service"

	"github.com/gophab/gophrame/module/system/domain"
	"github.com/gophab/gophrame/module/system/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SocialUserService struct {
	service.BaseService
	SocialUserRepository *repository.SocialUserRepository `inject:"socialUserRepository"`
}

var socialUserService *SocialUserService = &SocialUserService{}

func init() {
	eventbus.RegisterEventListener("USER_LOGIN", socialUserService.onUserLogin)
	eventbus.RegisterEventListener("SOCIAL_BIND_MOBILE", socialUserService.OnSocialBindUser)

	inject.InjectValue("socialUserService", socialUserService)

	security.RegisterPostHandlerFunc(socialUserService.ProcessSocialInfo)
}

func GetSocialUserService() *SocialUserService {
	return socialUserService
}

func (s *SocialUserService) GetById(id string) (*domain.SocialUser, error) {
	return s.SocialUserRepository.GetById(id)
}

func (s *SocialUserService) GetBySocialId(socialType string, socialId string) (*domain.SocialUser, error) {
	return s.SocialUserRepository.GetBySocialId(socialType, socialId)
}

func (s *SocialUserService) GetByUserId(socialType string, userId string) (*domain.SocialUser, error) {
	return s.SocialUserRepository.GetByUserId(socialType, userId)
}

func (s *SocialUserService) GetByMobile(socialType string, mobile string) (*domain.SocialUser, error) {
	return s.SocialUserRepository.GetByMobile(socialType, mobile)
}

func (s *SocialUserService) GetByEmail(socialType string, email string) (*domain.SocialUser, error) {
	return s.SocialUserRepository.GetByEmail(socialType, email)
}

func (s *SocialUserService) CreateSocialUser(socialUser *domain.SocialUser) (*domain.SocialUser, error) {
	socialUser.Status = util.IntAddr(consts.STATUS_VALID)
	if res := s.SocialUserRepository.Create(socialUser); res.Error == nil && res.RowsAffected > 0 {
		return socialUser, nil
	} else {
		return nil, res.Error
	}
}

func (s *SocialUserService) UpdateSocialUser(socialUser *domain.SocialUser) (*domain.SocialUser, error) {
	exists, err := s.GetById(socialUser.Id)
	if err != nil || exists == nil {
		return nil, err
	}

	updated := false

	if !util.StringPtrEquals(socialUser.Name, exists.Name) {
		exists.Name = socialUser.Name
		updated = true
	}
	if !util.StringPtrEquals(socialUser.NickName, exists.NickName) {
		exists.NickName = socialUser.NickName
		updated = true
	}
	if !util.StringPtrEquals(socialUser.Avatar, exists.Avatar) {
		exists.Avatar = socialUser.Avatar
		updated = true
	}
	if !util.StringPtrEquals(socialUser.Mobile, exists.Mobile) {
		exists.Mobile = socialUser.Mobile
		updated = true
	}
	if !util.StringPtrEquals(socialUser.Email, exists.Email) {
		exists.Email = socialUser.Email
		updated = true
	}
	if !util.StringPtrEquals(socialUser.Title, exists.Title) {
		exists.Title = socialUser.Title
		updated = true
	}
	if !util.StringPtrEquals(socialUser.Remark, exists.Remark) {
		exists.Remark = socialUser.Remark
		updated = true
	}
	if !util.StringPtrEquals(socialUser.OpenId, exists.OpenId) {
		exists.OpenId = socialUser.OpenId
		updated = true
	}
	if !util.StringPtrEquals(socialUser.SocialId, exists.SocialId) {
		exists.SocialId = socialUser.SocialId
		updated = true
	}
	if !util.StringPtrEquals(socialUser.UserId, exists.UserId) {
		exists.UserId = socialUser.UserId
		updated = true
	}
	if !util.PtrEquals(socialUser.Status, exists.Status) {
		exists.Status = socialUser.Status
		updated = true
	}

	if updated {
		if res := s.SocialUserRepository.Select("id").Omit("type", "login_times", "last_login_time", "last_login_ip", "created_time", "last_modified_time").Save(exists); res.Error == nil {
			return exists, nil
		} else {
			return nil, res.Error
		}
	}
	return exists, nil
}

func (s *SocialUserService) BoundSocialUser(socialUserId string, userId string, socialUser *domain.SocialUser) (*domain.SocialUser, error) {
	exists, err := s.GetById(socialUserId)
	if err != nil || exists == nil {
		return nil, err
	}

	exists.UserId = &userId

	var columns = core.M{
		"user_id": userId,
	}

	if socialUser != nil {
		if socialUser.Name != nil && !util.PtrEquals(exists.Name, socialUser.Name) {
			exists.Name = socialUser.Name
			columns["name"] = *socialUser.Name
		}

		if socialUser.NickName != nil && !util.PtrEquals(exists.NickName, socialUser.NickName) {
			exists.NickName = socialUser.NickName
			columns["nick_name"] = *socialUser.NickName
		}

		if socialUser.Avatar != nil && !util.PtrEquals(exists.Avatar, socialUser.Avatar) {
			exists.Avatar = socialUser.Avatar
			columns["avatar"] = *socialUser.Avatar
		}

		if socialUser.Mobile != nil && !util.PtrEquals(exists.Mobile, socialUser.Mobile) {
			exists.Mobile = socialUser.Mobile
			columns["mobile"] = *socialUser.Mobile
		}

		if socialUser.Email != nil && !util.PtrEquals(exists.Email, socialUser.Email) {
			exists.Email = socialUser.Email
			columns["email"] = *socialUser.Email
		}

		if socialUser.Title != nil && !util.PtrEquals(exists.Title, socialUser.Title) {
			exists.Title = socialUser.Title
			columns["title"] = *socialUser.Title
		}

		if socialUser.Remark != nil && !util.PtrEquals(exists.Remark, socialUser.Remark) {
			exists.Remark = socialUser.Remark
			columns["remark"] = *socialUser.Remark
		}
	}

	if err := s.SocialUserRepository.Changes(socialUserId, columns); err == nil {
		if exists.OpenId != nil && exists.Id != "sns:"+exists.Type+"_"+*exists.OpenId {
			s.BoundSocialUser("sns:"+exists.Type+"_"+*exists.OpenId, userId, socialUser)
		}
		return exists, nil
	} else {
		return nil, err
	}
}

func (s *SocialUserService) onUserLogin(event string, args ...any) {
	userId := ""
	data := map[string]string{}

	if len(args) > 0 {
		for i, v := range args {
			if i == 0 {
				userId = v.(string)
			}
			if i == 1 {
				data = v.(map[string]string)
			}
		}
	}

	if strings.HasPrefix(userId, "sns:") {
		if socialUser, err := s.GetById(userId); err != nil || socialUser == nil {
			return
		} else {
			if socialUser.UserId != nil {
				eventbus.PublishEvent("USER_LOGIN", *socialUser.UserId, data)
			}

			if _, err := s.SocialUserRepository.ConditionChanges(
				core.M{
					"id": socialUser.Id,
				},
				core.M{
					"last_login_time": time.Now(),
					"last_login_ip":   data["IP"],
					"login_times":     gorm.Expr("login_times+1"),
				},
			); err != nil {
				logger.Error("Save SocialUser Error: ", err.Error())
				return
			}
			// socialUser.LastLoginTime = util.TimeAddr(time.Now())
			// socialUser.LastLoginIp = util.StringAddr(data["IP"])
			// socialUser.LoginTimes = socialUser.LoginTimes + 1
			// if res := s.SocialUserRepository.Save(socialUser); res.Error != nil {
			// 	logger.Error("Save SocialUser Error: ", res.Error.Error())
			// 	return
			// }
		}
	}
}

func (s *SocialUserService) OnSocialBindUser(event string, args ...any) {
	var params = args[0].(core.M)

	if params != nil {
		social := params["social"]
		mobile := params["mobile"]
		email := params["email"]
		avatar := params["avatar"]
		nickName := params["nickName"]
		if userId, b := params["userId"]; b && userId != nil {
			var socialUser *domain.SocialUser

			if strings.HasPrefix(userId.(string), "sns:") {
				// 社交账号登录，更新相关信息
				socialUser, _ = s.GetById(userId.(string))
				if socialUser != nil {
					var columns = core.M{}
					if mobile != nil && (socialUser.Mobile == nil || *socialUser.Mobile != mobile.(string)) {
						columns["mobile"] = mobile.(string)
					}
					if email != nil && (socialUser.Email == nil || *socialUser.Email != email.(string)) {
						columns["email"] = email.(string)
					}
					if avatar != nil && (socialUser.Avatar == nil || *socialUser.Avatar != avatar.(string)) {
						columns["avatar"] = avatar.(string)
					}
					if nickName != nil && (socialUser.NickName == nil || *socialUser.NickName != nickName.(string)) {
						columns["nickName"] = nickName.(string)
					}
					if nickName != nil && socialUser.Name == nil {
						columns["name"] = nickName.(string)
					}

					if len(columns) > 0 {
						s.SocialUserRepository.Changes(userId.(string), columns)
					}
				}
			}

			if !strings.HasPrefix(userId.(string), "sns:") {
				// 用户登录，与社交账号里mobile为params["mobile"]进行绑定
				switch social.(string) {
				case "wxma", "wx", "wxmp":
					if mobile != nil {
						socialUser, _ = s.GetByMobile("wx", mobile.(string))
					} else if email != nil {
						socialUser, _ = s.GetByEmail("wx", email.(string))
					}
				}

				if socialUser != nil && socialUser.UserId == nil {
					su := &domain.SocialUser{}
					if avatar != nil {
						su.Avatar = util.StringAddr(avatar.(string))
					}
					if nickName != nil {
						su.NickName = util.StringAddr(nickName.(string))
					}

					s.BoundSocialUser(socialUser.Id, userId.(string), su)
				}
			}
		}
	}
}

func (s *SocialUserService) ProcessSocialInfo(context *gin.Context) {
	currentUserId := SecurityUtil.GetCurrentUserId(context)
	if strings.HasPrefix(currentUserId, "sns:") {
		// 1. social user
		if su, err := s.GetById(currentUserId); err == nil && su != nil {
			context.Set("open_id", util.NotNullString(su.OpenId))
		}
	} else {
		// 2. user
		if su, err := s.GetByUserId("wx", currentUserId); err == nil && su != nil {
			context.Set("open_id", util.NotNullString(su.OpenId))
		}
	}
}
