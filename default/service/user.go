package service

import (
	"errors"
	"strings"

	"github.com/gophab/gophrame/core/consts"
	"github.com/gophab/gophrame/core/eventbus"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/query"
	"github.com/gophab/gophrame/core/util"
	"github.com/gophab/gophrame/service"

	"github.com/gophab/gophrame/default/domain"
	"github.com/gophab/gophrame/default/repository"
	"github.com/gophab/gophrame/default/service/dto"

	"github.com/casbin/casbin/v2"
)

type UserService struct {
	service.BaseService
	UserRepository *repository.UserRepository `inject:"userRepository"`
	Enforcer       *casbin.SyncedEnforcer     `inject:"enforcer"`
}

var userService *UserService = &UserService{}

func GetUserService() *UserService {
	return userService
}

func init() {
	logger.Info("Initializing UserService...")
	logger.Debug("Inject UserService")
	inject.InjectValue("userService", userService)
	eventbus.RegisterEventListener("USER_LOGIN", userService.onUserLogin)
	logger.Info("Initialized UserService")
}

func (s *UserService) Check(username, password string) (bool, error) {
	return s.UserRepository.CheckUser(username, util.MD5(password))
}

func (s *UserService) CreateUser(user *dto.User) (*domain.User, error) {
	if user.Login != nil {
		if b, _ := s.UserRepository.CheckUserLogin(*user.Login); b {
			return nil, errors.New("用户名重复,请更改！")
		}
	}

	if user.Mobile != nil {
		if b, _ := s.UserRepository.CheckUserMobile(*user.Mobile); b {
			return nil, errors.New("手机号重复,请更改！")
		}
	}

	if user.Email != nil {
		if b, _ := s.UserRepository.CheckUserMobile(*user.Email); b {
			return nil, errors.New("邮箱重复,请更改！")
		}
	}

	if user.Status == nil {
		user.Status = util.IntAddr(consts.STATUS_VALID)
	}

	res, err := s.UserRepository.CreateUser(user.AsDomain())
	if err != nil {
		return nil, err
	}

	if user.InviteCode != nil {
		res.InviteCode = *user.InviteCode
	}

	eventbus.PublishEvent("USER_CREATED", res)
	return res, nil
}

func (s *UserService) UpdateUser(user *dto.User) (*domain.User, error) {
	if user.Id == nil {
		return nil, errors.New("Id为空")
	}

	exists, err := s.GetById(*user.Id)
	if err != nil || exists == nil {
		return nil, err
	}

	if user.Login != nil {
		b, _ := s.UserRepository.CheckUserLoginId(*user.Login, *user.Id)
		if b {
			return nil, errors.New("用户名重复,请更改！")
		}
		exists.Login = user.Login
	}

	if user.Mobile != nil {
		b, _ := s.UserRepository.CheckUserMobileId(*user.Mobile, *user.Id)
		if b {
			return nil, errors.New("手机号重复,请更改！")
		}
		exists.Mobile = user.Mobile
	}

	if user.Email != nil {
		b, _ := s.UserRepository.CheckUserEmailId(*user.Email, *user.Id)
		if b {
			return nil, errors.New("邮箱重复,请更改！")
		}
		exists.Email = user.Email
	}

	if user.Name != nil {
		exists.Name = user.Name
	}
	if user.Avatar != nil {
		exists.Avatar = user.Avatar
	}
	if user.Remark != nil {
		exists.Remark = user.Remark
	}

	if user.Status != nil {
		exists.Status = user.Status
	}

	if err = s.UserRepository.UpdateUser(exists); err == nil {
		eventbus.PublishEvent("USER_UPDATED", user)
	}

	return exists, err
}

func (s *UserService) Update(id string, column string, value interface{}) (*domain.User, error) {
	if res := s.UserRepository.Model(&domain.User{}).Where("id=?", id).Update(column, value); res.Error != nil {
		return nil, res.Error
	} else {
		if user, err := s.GetById(id); err == nil {
			eventbus.PublishEvent("USER_UPDATED", user)
			return user, err
		} else {
			return nil, err
		}
	}
}

func (s *UserService) Get(username string) (*domain.User, error) {
	user, err := s.UserRepository.GetUser(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetByLogin(login string) (*domain.User, error) {
	user, err := s.UserRepository.GetUserByLogin(login)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetByMobile(mobile string) (*domain.User, error) {
	user, err := s.UserRepository.GetUserByMobile(mobile)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetByEmail(email string) (*domain.User, error) {
	user, err := s.UserRepository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetById(id string) (*domain.User, error) {
	user, err := s.UserRepository.GetUserById(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetAll(user *dto.User, pageable query.Pageable) (int64, []domain.User) {
	if user.Id != nil {
		maps := make(map[string]interface{})
		maps["del_flag"] = false
		maps["id"] = user.Id
		return s.UserRepository.GetUsers(maps, pageable)
	} else {
		return s.UserRepository.GetUsers(user.GetMaps(), pageable)
	}
}

// 查询用户信息(带岗位)
func (a *UserService) GetAllWithOrganization(name string, pageable query.Pageable) (int64, []domain.UserWithOrganization) {
	return a.UserRepository.GetUserWithOrganizations(name, pageable)
}

func (s *UserService) Delete(id string) error {
	user, _ := s.GetById(id)
	if user != nil {
		err := s.UserRepository.DeleteUser(id)
		if err != nil {
			return err
		}

		if s.Enforcer != nil {
			s.Enforcer.DeleteUser(user.Id)
		}
	}
	return nil
}

func (s *UserService) ExistByID(id string) (bool, error) {
	return s.UserRepository.ExistUserByID(id)
}

func (s *UserService) Count(user *dto.User) (int64, error) {
	return s.UserRepository.GetUserTotal(user.GetMaps())
}

func (s *UserService) onUserLogin(event string, args ...interface{}) {
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

	if userId != "" && !strings.HasPrefix(userId, "sns:") {
		s.UserRepository.LogUserLogin(userId, data["IP"])
	}
}
