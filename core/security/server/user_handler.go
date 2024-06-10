package server

import (
	"context"

	SecurityModel "github.com/gophab/gophrame/core/security/model"
)

type IUserHandler interface {
	GetUserDetails(ctx context.Context, username string, password string) (*SecurityModel.UserDetails, error)
}

type IMobileUserHandler interface {
	GetMobileUserDetails(ctx context.Context, mobile string, code string) (*SecurityModel.UserDetails, error)
}

type IEmailUserHandler interface {
	GetEmailUserDetails(ctx context.Context, email string, code string) (*SecurityModel.UserDetails, error)
}

type ISocialUserHandler interface {
	GetSocialUserDetails(ctx context.Context, social string, code string) (*SecurityModel.UserDetails, error)
}
