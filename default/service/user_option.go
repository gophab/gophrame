package service

import (
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/service"

	"github.com/gophab/gophrame/default/domain"
	"github.com/gophab/gophrame/default/repository"
)

type UserOptionService struct {
	service.BaseService
	UserOptionRepository *repository.UserOptionRepository `inject:"userOptionRepository"`
}

var userOptionService = &UserOptionService{}

func init() {
	inject.InjectValue("userOptionService", userOptionService)
}

var defaultUserOptions = map[string]string{}

func (s *UserOptionService) GetDefaultUserOptions() (*domain.UserOptions, error) {
	result := &domain.UserOptions{UserId: "DEFAULT", Options: make(map[string]domain.UserOption)}

	// DEFAULT
	for k, v := range defaultUserOptions {
		result.Options[k] = domain.UserOption{
			UserId: "DEFAULT",
			Option: domain.Option{
				Name:      k,
				Value:     v,
				ValueType: "STRING",
			},
		}
	}

	// DEFAULT in DB
	resultDB, err := s.UserOptionRepository.GetUserOptions("DEFAULT")
	if err != nil {
		return nil, err
	}

	if resultDB != nil && len(resultDB.Options) > 0 {
		for k, v := range resultDB.Options {
			v.UserId = "DEFAULT"
			result.Options[k] = v
		}
	}

	return result, nil
}

func (s *UserOptionService) GetUserOptions(userId string) (*domain.UserOptions, error) {
	result, err := s.GetDefaultUserOptions()
	if err != nil {
		return nil, err
	}

	resultDB, err := s.UserOptionRepository.GetUserOptions(userId)
	if err != nil {
		return nil, err
	}

	result.UserId = userId

	if resultDB != nil && len(resultDB.Options) > 0 {
		for k, v := range resultDB.Options {
			v.UserId = userId
			result.Options[k] = v
		}
		return result, nil
	}

	return result, nil
}

func (s *UserOptionService) AddUserOption(option *domain.UserOption) (*domain.UserOption, error) {
	if res := s.UserOptionRepository.Save(option); res.Error == nil && res.RowsAffected > 0 {
		return option, nil
	} else {
		return nil, res.Error
	}
}

func (s *UserOptionService) AddUserOptions(options []domain.UserOption) (*[]domain.UserOption, error) {
	var result = make([]domain.UserOption, len(options))
	for i, option := range options {
		if res := s.UserOptionRepository.Save(option); res.Error != nil {
			return nil, res.Error
		}
		result[i] = option
	}
	return &result, nil
}

func (s *UserOptionService) RemoveAllUserOptions(userId string) error {
	return s.UserOptionRepository.RemoveAllUserOptions(userId)
}

func (s *UserOptionService) RemoveUserOption(userId string, key string) (*domain.UserOption, error) {
	return nil, s.UserOptionRepository.Delete(&domain.UserOption{UserId: userId, Option: domain.Option{Name: key}}).Error
}

func (s *UserOptionService) SetUserOption(userId string, key string, value string) (*domain.UserOption, error) {
	var option = domain.UserOption{
		UserId: userId,
		Option: domain.Option{
			Name:      key,
			Value:     value,
			ValueType: "STRING",
		},
	}

	if res := s.UserOptionRepository.Save(&option); res.Error == nil && res.RowsAffected > 0 {
		return &option, nil
	} else {
		return nil, res.Error
	}
}

func (s *UserOptionService) SetUserOptions(userOptions *domain.UserOptions) (*domain.UserOptions, error) {
	// 1. Remove User Options
	if err := s.RemoveAllUserOptions(userOptions.UserId); err != nil {
		return nil, err
	}

	// 2. Save
	var options []domain.UserOption
	for _, v := range userOptions.Options {
		v.UserId = userOptions.UserId
		options = append(options, v)
	}

	if _, err := s.AddUserOptions(options); err != nil {
		return nil, err
	}

	return userOptions, nil
}
