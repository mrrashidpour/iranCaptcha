package iranCaptcha

import (
	"image/color"
	"math"
	"math/rand/v2"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

const (
	//TxtNumbers chacters for numbers.
	TxtNumbers = "0123456789"
	//TxtAlphabet characters for alphabet.
	TxtAlphabet = "ABCDEFGHJKMNOQRSTUVXYZabcdefghjkmnoqrstuvxyz"
)

// idCounter برای تولید ID سریع و یکتا
var idCounter atomic.Uint64

// RandText creates random text of given size.
func RandText(size int, sourceChars string) string {
	if sourceChars == "" || size == 0 {
		return ""
	}

	if size >= len(sourceChars) {
		sourceChars = strings.Repeat(sourceChars, size)
	}

	sourceRunes := []rune(sourceChars)
	sourceLength := len(sourceRunes)

	text := make([]rune, size)
	for i := range text {
		text[i] = sourceRunes[rand.IntN(sourceLength)]
	}
	return string(text)
}

// random get random number between min and max.
func random(min int64, max int64) float64 {
	return float64(min) + rand.Float64()*float64(max-min)
}

// RandDeepColor get random deep color.
func RandDeepColor() color.RGBA {
	randColor := RandColor()

	increase := float64(30 + rand.IntN(255))
	red := math.Abs(math.Min(float64(randColor.R)-increase, 255))
	green := math.Abs(math.Min(float64(randColor.G)-increase, 255))
	blue := math.Abs(math.Min(float64(randColor.B)-increase, 255))

	return color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(255)}
}

// RandLightColor get random light color.
func RandLightColor() color.RGBA {
	red := rand.IntN(55) + 200
	green := rand.IntN(55) + 200
	blue := rand.IntN(55) + 200
	return color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(255)}
}

// RandDarkColor get random dark color.
func RandDarkColor() color.RGBA {
	maxColor := 130
	red := rand.IntN(maxColor)
	green := rand.IntN(maxColor)
	blue := rand.IntN(maxColor)
	return color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(255)}
}

// RandColor get random color.
func RandColor() color.RGBA {
	red := rand.IntN(255)
	green := rand.IntN(255)
	var blue int
	if (red + green) > 400 {
		blue = 0
	} else {
		blue = 400 - green - red
	}
	if blue > 255 {
		blue = 255
	}
	return color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(255)}
}

func randIntRange(from, to int) int {
	// rand.IntN panics if n <= 0.
	if to-from <= 0 {
		return from
	}
	return rand.IntN(to-from) + from
}

func randFloat64Range(from, to float64) float64 {
	return rand.Float64()*(to-from) + from
}

// RandomId یک ID سریع و یکتا برمی‌گرداند (counter + timestamp).
// فقط برای یکتایی، نه رمزنگاری‌امن.
func RandomId() string {
	n := idCounter.Add(1)
	return strconv.FormatUint(uint64(time.Now().UnixNano())<<16|n, 36)
}
