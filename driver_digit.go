package iranCaptcha

import (
	"image/color"
	"math/rand/v2"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/mojocn/base64Captcha"
	"github.com/mrrashidpour/iranCaptcha/fonts"
)

type DriverDigit struct {
	Height       int
	Width        int
	NoiseCount   int
	Length       int
	Source       string
	BgColor      *color.RGBA
	fontsStorage base64Captcha.FontsStorage
	Fonts        []string
	fontsArray   []*truetype.Font
}

func NewDriverDigit(height int, width int, noiseCount int, length int, bgColor *color.RGBA, font string) *DriverDigit {
	fontsStorage := DefaultEmbeddedFonts
	if font == "" {
		font = fonts.FaNum
	}
	if bgColor == nil {
		bgColor = &color.RGBA{R: 0, G: 0, B: 0, A: 0}
	}

	tf := fontsStorage.LoadFontByName("fonts/" + font)
	tfs := []*truetype.Font{tf}

	return &DriverDigit{
		Height:       height,
		Width:        width,
		NoiseCount:   noiseCount,
		Length:       length,
		Source:       "0123456789",
		BgColor:      bgColor,
		fontsStorage: fontsStorage,
		fontsArray:   tfs,
		Fonts:        []string{font},
	}
}

func (d *DriverDigit) ConvertFonts() *DriverDigit {
	if d.fontsStorage == nil {
		d.fontsStorage = base64Captcha.DefaultEmbeddedFonts
	}
	tfs := []*truetype.Font{}
	for _, fff := range d.Fonts {
		tf := d.fontsStorage.LoadFontByName("fonts/" + fff)
		tfs = append(tfs, tf)
	}
	d.fontsArray = tfs
	return d
}

func (d *DriverDigit) GenerateIdQuestionAnswer() (id, content, answer string) {
	id = RandomId()
	content = RandText(d.Length, d.Source)
	return id, content, content
}

func (d *DriverDigit) DrawCaptcha(content string) (item base64Captcha.Item, err error) {
	var bgc color.RGBA
	if d.BgColor != nil {
		bgc = *d.BgColor
	} else {
		bgc = RandLightColor()
	}
	itemChar := NewItemChar(d.Width, d.Height, bgc)

	itemChar.drawSlimLine(rand.IntN(3))

	if d.NoiseCount > 0 {
		source := ".,"
		noise := RandText(d.NoiseCount, strings.Repeat(source, d.NoiseCount))
		err = itemChar.drawNoise(noise, d.fontsArray)
		if err != nil {
			return
		}
	}

	err = itemChar.drawText(content, d.fontsArray)
	if err != nil {
		return
	}

	return itemChar, nil
}
