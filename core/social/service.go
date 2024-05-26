package social

import "context"

type SocialService interface {
	GetSocialUserByCode(ctx context.Context, socialChannelId string, code string) *SocialUser
}
