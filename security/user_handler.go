package security

import (
	"context"

	"github.com/gophab/gophrame/core/inject"

	SecurityModel "github.com/gophab/gophrame/core/security/model"
)

type UserHandler interface {
	GetUserDetails(ctx context.Context, username, password string) (*SecurityModel.UserDetails, error)
	GetMobileUserDetails(ctx context.Context, mobile string, code string) (*SecurityModel.UserDetails, error)
	GetEmailUserDetails(ctx context.Context, email string, code string) (*SecurityModel.UserDetails, error)
	GetSocialUserDetails(ctx context.Context, socialChannelId string, code string) (*SecurityModel.UserDetails, error)
}

func RegisterUserHandler(handler UserHandler) {
	inject.InjectValue("userHandler", handler)
}
