package SecurityModel

import "github.com/gophab/gophrame/core/util/collection"

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

func (u *UserDetails) HasRole(role string) bool {
	if len(u.Roles) > 0 {
		b, _ := collection.Contains(u.Roles, role)
		return b
	}

	return false
}
