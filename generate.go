package iranCaptcha

import (
	"github.com/mojocn/base64Captcha"
)

func Generate(driver base64Captcha.Driver, store base64Captcha.Store) (id, b64s, answer string, err error) {

	newCaptcha := base64Captcha.NewCaptcha(driver, store)

	return newCaptcha.Generate()

}
