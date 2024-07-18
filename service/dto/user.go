package dto

type User struct {
	Id            *string `json:"id"`
	Login         *string `json:"login"`
	Mobile        *string `json:"mobile"`
	Email         *string `json:"email"`
	Password      *string `json:"-"`
	PlainPassword *string `json:"password,omitempty"`
	Name          *string `json:"name"`
	InviteCode    *string `json:"inviteCode"`
	InviterId     *string `json:"inviterId"`
	SocialId      *string `json:"socialId"`
	Admin         *bool   `json:"admin"`
	TenantId      *string `json:"tenantId"`
}
