package domain

type UserOption struct {
	Option
	UserId string `gorm:"column:user_id"`
}

func (*UserOption) TableName() string {
	return "sys_user_option"
}

type UserOptions struct {
	UserId  string
	Options map[string]UserOption
}

func (s *UserOptions) GetOption(name string) (string, bool) {
	if option, b := s.Options[name]; b {
		return option.Value, true
	} else {
		return "", false
	}
}

func (s *UserOptions) SetOption(name string, value string) bool {
	if option, b := s.Options[name]; b {
		option.Value = value
		s.Options[name] = option
	} else {
		s.Options[name] = UserOption{
			UserId: s.UserId,
			Option: Option{
				Name:      name,
				Value:     value,
				ValueType: "STRING",
			},
		}
	}
	return true
}
