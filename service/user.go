package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/wjshen/gophrame/core/eventbus"
	"github.com/wjshen/gophrame/core/inject"
	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/query"
	"github.com/wjshen/gophrame/core/util"

	"github.com/wjshen/gophrame/domain"
	"github.com/wjshen/gophrame/repository"
	"github.com/wjshen/gophrame/service/dto"

	"github.com/casbin/casbin/v2"
)

type UserService struct {
	BaseService
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
		user.Status = &domain.STATUS_VALID
	}

	res, err := s.UserRepository.CreateUser(&user.User)
	if err != nil {
		return nil, err
	}

	if user.InviteCode != "" {
		res.InviteCode = user.InviteCode
	}

	eventbus.PublishEvent("USER_REGISTERED", res)
	return res, nil
}

func (s *UserService) UpdateUser(user *dto.User) (*domain.User, error) {
	exists, err := s.GetById(user.Id)
	if err != nil || exists == nil {
		return nil, err
	}

	if user.Login != nil {
		b, _ := s.UserRepository.CheckUserLoginId(*user.Login, user.Id)
		if b {
			return nil, errors.New("用户名重复,请更改！")
		}
		exists.Login = user.Login
	}

	if user.Mobile != nil {
		b, _ := s.UserRepository.CheckUserMobileId(*user.Mobile, user.Id)
		if b {
			return nil, errors.New("手机号重复,请更改！")
		}
		exists.Mobile = user.Mobile
	}

	if user.Email != nil {
		b, _ := s.UserRepository.CheckUserEmailId(*user.Email, user.Id)
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
		err = s.LoadPolicy(user.Id)
	}

	return exists, err
}

func (s *UserService) Update(id string, column string, value interface{}) (*domain.User, error) {
	if res := s.UserRepository.Model(&domain.User{}).Where("id=?", id).Update(column, value); res.Error != nil {
		return nil, res.Error
	} else {
		return s.GetById(id)
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
	if user.Id != "" {
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

// LoadAllPolicy 加载所有的用户策略
func (s *UserService) LoadAllPolicy() error {
	if s.Enforcer != nil {
		users, err := s.UserRepository.GetUsersAll()
		if err != nil {
			return err
		}
		for _, user := range users {
			if len(user.Roles) != 0 {
				err = s.LoadPolicy(user.Id)
				if err != nil {
					return err
				}
			}
		}
		fmt.Println("角色权限关系", s.Enforcer.GetGroupingPolicy())
	}
	return nil
}

// LoadPolicy 加载用户权限策略
func (s *UserService) LoadPolicy(id string) error {
	if s.Enforcer != nil {
		user, err := s.UserRepository.GetUserById(id)
		if err != nil {
			return err
		}

		s.Enforcer.DeleteRolesForUser(user.Id)
		for _, ro := range user.Roles {
			s.Enforcer.AddRoleForUser(user.Id, ro.Name)
		}
		fmt.Println("更新角色权限关系", s.Enforcer.GetGroupingPolicy())
	}
	return nil
}

func (s *UserService) onUserLogin(args ...interface{}) {
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
