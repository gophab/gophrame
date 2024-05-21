package server

import (
	"context"

	"github.com/wjshen/gophrame/core/security/model"
)

type IUserHandler interface {
	GetUserDetails(ctx context.Context, username string, password string) (*model.UserDetails, error)
}

type IMobileUserHandler interface {
	GetMobileUserDetails(ctx context.Context, mobile string, code string) (*model.UserDetails, error)
}

type IEmailUserHandler interface {
	GetEmailUserDetails(ctx context.Context, email string, code string) (*model.UserDetails, error)
}

type ISocialUserHandler interface {
	GetSocialUserDetails(ctx context.Context, social string, code string) (*model.UserDetails, error)
}
