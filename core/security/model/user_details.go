package SecurityModel

type UserDetails struct {
	UserId   *string
	Name     *string
	Avatar   *string
	Login    *string
	Mobile   *string
	Email    *string
	SocialId *string
	TenantId *string
	Admin    bool
	Roles    []string
}
