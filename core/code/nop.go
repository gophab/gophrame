package code

type nop struct{}

func (*nop) CreateRequest(id string) error {
	return nil
}

func (*nop) CreateCode(id string, scene string, code string) error {
	return nil
}

func (*nop) GetCode(id string, scene string, remove bool) (string, bool) {
	return "", false
}

func (*nop) RemoveCode(id string, scene string) {
}

func (*nop) SendVerificationCode(dest string, scene string, code string) error {
	return nil
}

var Nop = nop{}
