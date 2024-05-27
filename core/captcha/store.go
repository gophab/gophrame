package captcha

import (
	"github.com/dchest/captcha"
	"github.com/gophab/gophrame/core/code"
	"github.com/mojocn/base64Captcha"
)

type CaptchaStoreAdpter struct {
	captcha.Store
	CodeStore code.CodeStore
}

// Set sets the digits for the captcha id.
func (a *CaptchaStoreAdpter) Set(id string, digits []byte) {
	a.CodeStore.CreateCode(id, "", string(digits[:]))
}

// Get returns stored digits for the captcha id. Clear indicates
// whether the captcha must be deleted from the store.
func (a *CaptchaStoreAdpter) Get(id string, clear bool) (digits []byte) {
	if v, b := a.CodeStore.GetCode(id, "", clear); b {
		return []byte(v)
	} else {
		return nil
	}
}

type Base64CaptchaStoreAdapter struct {
	base64Captcha.Store
	CodeStore code.CodeStore
}

// Set sets the digits for the captcha id.
func (a *Base64CaptchaStoreAdapter) Set(id string, value string) error {
	return a.CodeStore.CreateCode(id, "", value)
}

// Get returns stored digits for the captcha id. Clear indicates
// whether the captcha must be deleted from the store.
func (a *Base64CaptchaStoreAdapter) Get(id string, clear bool) string {
	if v, b := a.CodeStore.GetCode(id, "", clear); b {
		return v
	} else {
		return ""
	}
}

// Verify captcha's answer directly
func (a *Base64CaptchaStoreAdapter) Verify(id, answer string, clear bool) bool {
	if v, b := a.CodeStore.GetCode(id, "", clear); b {
		return v == answer
	} else {
		return false
	}
}
