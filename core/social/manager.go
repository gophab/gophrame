package social

import "github.com/gophab/gophrame/core/inject"

type SocialManager struct {
	Services map[string]SocialService
}

func GetSocialService(social string) SocialService {
	return socialManager.GetSocialService(social)
}

func RegisterSocialService(social string, service SocialService) error {
	return socialManager.RegisterSocialService(social, service)
}

var socialManager = &SocialManager{
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

func init() {
	inject.InjectValue("socialManager", socialManager)
}
