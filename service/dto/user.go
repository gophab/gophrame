package dto

type User struct {
	Id            *string `json:"id"`
	Login         *string `json:"login,omitempty"`
	Mobile        *string `json:"mobile,omitempty"`
	Email         *string `json:"email,omitempty"`
	Password      *string `json:"-"`
	PlainPassword *string `json:"password,omitempty"`
	Name          *string `json:"name,omitempty"`
	InviteCode    *string `json:"inviteCode,omitempty"`
	InviterId     *string `json:"inviterId,omitempty"`
	SocialId      *string `json:"socialId,omitempty"`
	Admin         *bool   `json:"admin,omitempty"`
	TenantId      *string `json:"tenantId,omitempty"`
}
