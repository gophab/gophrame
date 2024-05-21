package code

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type CodeStore interface {
	CreateRequest(phone string) error
	CreateCode(mobile string, scene string, code string) error
	GetCode(mobile string, scene string) (string, bool)
	RemoveCode(mobile string, scene string)
}

type CodeSender interface {
	SendVerificationCode(dest string, scene string, code string) error
}

type nop struct{}

func (*nop) CreateRequest(phone string) error {
	return nil
}

func (*nop) CreateCode(mobile string, scene string, code string) error {
	return nil
}

func (*nop) GetCode(mobile string, scene string) (string, bool) {
	return "", false
}

func (*nop) RemoveCode(mobile string, scene string) {
}

func (*nop) SendVerificationCode(dest string, scene string, code string) error {
	return nil
}

var Nop = nop{}

type CodeValidator interface {
	GetStore() CodeStore
	GetSender() CodeSender
	GenerateCode(target CodeValidator, dest string, scene string) (string, error)
	GetVerificationCode(target CodeValidator, dest string, scene string) (string, bool)
	CheckCode(target CodeValidator, dest string, scene string, code string) bool
}

func NewValidator(sender CodeSender, store CodeStore) CodeValidator {
	return &Validator{Sender: sender, Store: store}
}

type Validator struct {
	Store  CodeStore
	Sender CodeSender
}

func (v *Validator) GetStore() CodeStore {
	if v.Store != nil {
		return v.Store
	}
	return &Nop
}

func (v *Validator) GetSender() CodeSender {
	if v.Sender != nil {
		return v.Sender
	}
	return &Nop
}

func (v *Validator) GenerateCode(target CodeValidator, dest string, scene string) (string, error) {
	if err := target.GetStore().CreateRequest(dest); err != nil {
		return "", errors.New("请稍候重试")
	}

	// 新生成
	code := v.CreateRandCode()

	// 保存至缓存
	if err := target.GetStore().CreateCode(dest, scene, code); err != nil {
		return "", err
	}

	// 使用发送器
	if err := target.GetSender().SendVerificationCode(dest, scene, code); err != nil {
		return "", err
	}

	return code, nil
}

// 创建6位随机数
func (v *Validator) CreateRandCode() string {
	return fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
}

func (v *Validator) GetVerificationCode(target CodeValidator, dest string, scene string) (string, bool) {
	return target.GetStore().GetCode(dest, scene)
}

func (v *Validator) CheckCode(target CodeValidator, dest string, scene string, code string) bool {
	if cached, b := v.GetVerificationCode(target, dest, scene); b {
		return cached == code
	}
	return false
}
