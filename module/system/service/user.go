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
	CommonDTO "github.com/gophab/gophrame/service/dto"

	"github.com/gophab/gophrame/module/system/domain"
	"github.com/gophab/gophrame/module/system/repository"
	"github.com/gophab/gophrame/module/system/service/dto"

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
	inject.InjectValue("commonUserService", commonUserService)
	_ = (service.UserService)(commonUserService)
	eventbus.RegisterEventListener("USER_LOGIN", userService.onUserLogin)
	logger.Info("Initialized UserService")
}

func (s *UserService) GetById(id string) (*domain.User, error) {
	user, err := s.UserRepository.GetUserById(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetByIds(ids []string) ([]*domain.User, error) {
	users, err := s.UserRepository.GetUserByIds(ids)
	if err != nil {
		return nil, err
	}

	return users, nil
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

func (s *UserService) Create(user *dto.User) (*domain.User, error) {
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

	if user.PlainPassword != nil && (user.Password == nil || *user.Password == "") {
		user.Password = user.PlainPassword
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

func (s *UserService) Update(user *dto.User) (*domain.User, error) {
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

	if user.Password != nil {
		exists.Password = *user.Password
	}

	if err = s.UserRepository.UpdateUser(exists); err == nil {
		eventbus.PublishEvent("USER_UPDATED", user)
	}

	return exists, err
}

func (s *UserService) Patch(id string, column string, value interface{}) (*domain.User, error) {
	if column == "password" && value != nil {
		value = util.SHA1(value.(string))
	}
	if res := s.UserRepository.Model(&domain.User{}).Where("id=?", id).UpdateColumn(util.DbFieldName(column), value); res.Error != nil {
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

func (s *UserService) PatchAll(id string, kv map[string]interface{}) (*domain.User, error) {
	if kv["password"] != nil {
		kv["password"] = util.SHA1(kv["password"].(string))
	}

	kv["id"] = id
	if res := s.UserRepository.Model(&domain.User{}).Where("id=?", id).UpdateColumns(util.DbFields(kv)); res.Error != nil {
		return nil, res.Error
	}

	if user, err := s.GetById(id); err != nil {
		return nil, err
	} else {
		eventbus.PublishEvent("USER_UPDATED", user)
		return user, err
	}
}

func (s *UserService) Delete(user *domain.User) error {
	err := s.UserRepository.DeleteUser(user.Id)
	if err != nil {
		return err
	}

	if s.Enforcer != nil {
		s.Enforcer.DeleteUser(user.Id)
	}
	return nil
}

func (s *UserService) DeleteById(id string) error {
	user, _ := s.GetById(id)
	if user != nil {
		return s.Delete(user)
	}
	return nil
}

func (s *UserService) Check(username, password string) (bool, error) {
	return s.UserRepository.CheckUser(username, util.MD5(password))
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

func (a *UserService) Find(conds map[string]interface{}, pageable query.Pageable) (int64, []*domain.User) {
	return a.UserRepository.Find(util.DbFields(conds), pageable)
}

// 查询用户信息(带岗位)
func (a *UserService) GetAllWithOrganization(name string, pageable query.Pageable) (int64, []domain.UserWithOrganization) {
	return a.UserRepository.GetUserWithOrganizations(name, pageable)
}

func (s *UserService) ExistByID(id string) (bool, error) {
	return s.UserRepository.ExistUserByID(id)
}

func (s *UserService) Count(user *dto.User) (int64, error) {
	return s.UserRepository.GetUserTotal(user.GetMaps())
}

func (s *UserService) ResetUserPassword(id string) (bool, error) {
	res := s.UserRepository.Model(&domain.User{}).
		Where("id=?", id).
		UpdateColumn("password", util.SHA1("123456"))
	if res.Error != nil {
		return false, res.Error
	}
	return true, nil
}

func (s *UserService) ChangeUserPassword(id string, oldpassword, password string) (bool, error) {
	res := s.UserRepository.Model(&domain.User{}).
		Where("id=?", id).
		Where("password=?", util.SHA1(oldpassword)).
		UpdateColumn("password", util.SHA1(password))
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
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

type CommonUserService struct{}

var commonUserService = &CommonUserService{}

func (s *CommonUserService) CreateUser(user *CommonDTO.User) (*CommonDTO.User, error) {
	var dtoUser = dto.User{
		User: *user,
	}
	if result, err := userService.Create(&dtoUser); err != nil {
		return nil, err
	} else {
		user.Id = util.StringAddr(result.Id)
		return user, nil
	}
}

func (s *CommonUserService) GetById(id string) (*CommonDTO.User, error) {
	if result, err := userService.GetById(id); err != nil {
		return nil, err
	} else if result != nil {
		var user = &CommonDTO.User{
			Id:         &result.Id,
			InviteCode: &result.InviteCode,
			InviterId:  result.InviterId,
			Name:       result.Name,
			Mobile:     result.Mobile,
			Email:      result.Email,
			Admin:      &result.Admin,
			TenantId:   &result.TenantId,
		}
		return user, nil
	} else {
		return nil, nil
	}
}

func (s *CommonUserService) GetByIds(ids []string) ([]*CommonDTO.User, error) {
	_, list := userService.Find(map[string]interface{}{
		"ids": ids,
	}, nil)

	if len(list) > 0 {
		var results = make([]*CommonDTO.User, 0)
		for _, result := range list {
			var user = &CommonDTO.User{
				Id:         &result.Id,
				InviteCode: &result.InviteCode,
				InviterId:  result.InviterId,
				Name:       result.Name,
				Mobile:     result.Mobile,
				Email:      result.Email,
				Admin:      &result.Admin,
				TenantId:   &result.TenantId,
			}
			results = append(results, user)
		}
		return results, nil
	}

	return []*CommonDTO.User{}, nil
}

func (s *CommonUserService) GetByMobile(mobile string) (*CommonDTO.User, error) {
	if result, err := userService.GetByMobile(mobile); err != nil {
		return nil, err
	} else if result != nil {
		var user = &CommonDTO.User{
			Id:         &result.Id,
			InviteCode: &result.InviteCode,
			InviterId:  result.InviterId,
			Name:       result.Name,
			Mobile:     result.Mobile,
			Email:      result.Email,
			Admin:      &result.Admin,
			TenantId:   &result.TenantId,
		}
		return user, nil
	} else {
		return nil, nil
	}
}

func (s *CommonUserService) GetByEmail(email string) (*CommonDTO.User, error) {
	if result, err := userService.GetByEmail(email); err != nil {
		return nil, err
	} else if result != nil {
		var user = &CommonDTO.User{
			Id:         &result.Id,
			InviteCode: &result.InviteCode,
			InviterId:  result.InviterId,
			Name:       result.Name,
			Mobile:     result.Mobile,
			Email:      result.Email,
			Admin:      &result.Admin,
			TenantId:   &result.TenantId,
		}
		return user, nil
	} else {
		return nil, nil
	}
}
