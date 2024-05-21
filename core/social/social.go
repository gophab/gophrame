package social

import (
	"context"

	"github.com/wjshen/gophrame/domain"
)

type SocialService interface {
	GetSocialUserByCode(ctx context.Context, socialChannelId string, code string) *domain.SocialUser
}

type SocialManager struct {
	Services map[string]SocialService
}

func GetSocialService(social string) SocialService {
	return Manager.GetSocialService(social)
}

func RegisterSocialService(social string, service SocialService) error {
	return Manager.RegisterSocialService(social, service)
}

var Manager = &SocialManager{
	Services: make(map[string]SocialService, 10),
}

func (m *SocialManager) GetSocialService(social string) SocialService {
	if service, b := m.Services[social]; b {
		return service
	}
	return nil
}

func (m *SocialManager) RegisterSocialService(social string, service SocialService) error {
	m.Services[social] = service
	return nil
}
