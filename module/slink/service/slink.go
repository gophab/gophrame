package service

import (
	"strings"
	"time"

	"github.com/gophab/gophrame/module/slink/config"
	"github.com/gophab/gophrame/module/slink/domain"
	"github.com/gophab/gophrame/module/slink/repository"

	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/starter"
	"github.com/gophab/gophrame/core/util"
	"github.com/gophab/gophrame/service"
)

type ShortLinkService struct {
	service.BaseService
	ShortLinkRepository *repository.ShortLinkRepository `inject:"shortLinkRepository"`
}

var shortLinkService = &ShortLinkService{}

func init() {
	inject.InjectValue("shortLinkService", shortLinkService)
	starter.RegisterStarter(shortLinkService.Start)
}

func (s *ShortLinkService) GetById(id string) (*domain.ShortLink, error) {
	return s.ShortLinkRepository.GetById(id)
}

func (s *ShortLinkService) DeleteById(id string) error {
	return s.ShortLinkRepository.DeleteById(id)
}

func (s *ShortLinkService) GetByKey(key string) (*domain.ShortLink, error) {
	return s.ShortLinkRepository.GetByKey(key)
}

func (s *ShortLinkService) DeleteByKey(key string) error {
	return s.ShortLinkRepository.DeleteByKey(key)
}

func (s *ShortLinkService) Generate(url string) *domain.ShortLink {
	result, _ := s.GenerateShortLink("", url, config.Setting.Expired)
	if result != nil {
		return result
	} else {
		return &domain.ShortLink{Url: url, FullPath: url}
	}
}

func (s *ShortLinkService) GenerateShortLink(name string, url string, duration time.Duration) (*domain.ShortLink, error) {
	var result = domain.ShortLink{
		Name:        name,
		Key:         util.GenerateRandomString(5),
		Url:         url,
		ExpiredTime: util.TimeAddr(time.Now().Add(duration)),
	}

	root, _ := strings.CutSuffix(config.Setting.BaseUrl, "/")
	result.FullPath = root + "/"
	if config.Setting.Context != "" {
		result.FullPath = result.FullPath + config.Setting.Context + "/"
	}
	result.FullPath = result.FullPath + result.Key

	return s.ShortLinkRepository.CreateShortLink(&result)
}

func (s *ShortLinkService) Start() {
	logger.Info("Starting shortlink expire routine...")
	go func() {
		for {
			s.ShortLinkRepository.ExpireExpiredShortLinks()
			time.Sleep(time.Hour * time.Duration(1))
		}
	}()
}
