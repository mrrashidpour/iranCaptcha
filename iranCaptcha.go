package iranCaptcha

import (
	"image/color"
	"time"

	"github.com/mojocn/base64Captcha"
	"github.com/mrrashidpour/iranCaptcha/services"
)

func NewMemoryStore(collectNum int, expiration time.Duration) base64Captcha.Store {
	return base64Captcha.NewMemoryStore(collectNum, expiration)
}

func NewDriverDigit2(height int, width int, noiseCount int, length int, bgColor *color.RGBA, font string) *services.DriverDigit2 {
	return services.NewDriverDigit2(height, width, noiseCount, length, bgColor, font)
}

func Generate(driver base64Captcha.Driver, store base64Captcha.Store) (id, b64s, answer string, err error) {

	newCaptcha := base64Captcha.NewCaptcha(driver, store)

	return newCaptcha.Generate()

}
