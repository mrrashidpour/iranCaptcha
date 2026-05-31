package services

import (
	"image/color"
	"math/rand"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/mojocn/base64Captcha"
	"github.com/mrrashidpour/iranCaptcha/config"
	"github.com/mrrashidpour/iranCaptcha/fonts"
)

type DriverDigit2 struct {
	// Height png height in pixel.
	Height int

	// Width Captcha png width in pixel.
	Width int

	//NoiseCount text noise count.
	NoiseCount int

	//Length random string length.
	Length int

	//Source is a unicode which is the rand string from.
	Source string

	//BgColor captcha image background color (optional)
	BgColor *color.RGBA

	//fontsStorage font storage (optional)
	fontsStorage base64Captcha.FontsStorage

	//Fonts loads by name see fonts.go's comment
	Fonts      []string
	fontsArray []*truetype.Font
}

// NewDriverDigit2 creates driver
func NewDriverDigit2(height int, width int, noiseCount int, length int, bgColor *color.RGBA, font string) *DriverDigit2 {

	fontsStorage := config.DefaultEmbeddedFonts

	if font == "" {
		font = fonts.BTitrBd
	}

	if bgColor == nil {
		bgColor = &color.RGBA{R: 0, G: 0, B: 0, A: 0}
	}

	tf := fontsStorage.LoadFontByName("fonts/" + font)

	tfs := []*truetype.Font{tf}

	return &DriverDigit2{
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

// ConvertFonts loads fonts by names
func (d *DriverDigit2) ConvertFonts() *DriverDigit2 {
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

// GenerateIdQuestionAnswer creates id,content and answer
func (d *DriverDigit2) GenerateIdQuestionAnswer() (id, content, answer string) {
	id = RandomId()
	content = RandText(d.Length, d.Source)
	return id, content, content
}

// DrawCaptcha draws captcha item
func (d *DriverDigit2) DrawCaptcha(content string) (item base64Captcha.Item, err error) {

	var bgc color.RGBA
	if d.BgColor != nil {
		bgc = *d.BgColor
	} else {
		bgc = RandLightColor()
	}
	itemChar := NewItemChar(d.Width, d.Height, bgc)

	itemChar.drawSlimLine(rand.Intn(3))

	//draw noise
	if d.NoiseCount > 0 {
		source := ".,"
		noise := RandText(d.NoiseCount, strings.Repeat(source, d.NoiseCount))
		err = itemChar.drawNoise(noise, d.fontsArray)
		if err != nil {
			return
		}
	}

	//draw content
	err = itemChar.drawText(content, d.fontsArray)
	if err != nil {
		return
	}

	return itemChar, nil
}
