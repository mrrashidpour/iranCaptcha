package iranCaptcha

import (
	"time"

	"github.com/mojocn/base64Captcha"
)

func NewMemoryStore(collectNum int, expiration time.Duration) base64Captcha.Store {
	return base64Captcha.NewMemoryStore(collectNum, expiration)
}
