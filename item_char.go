package iranCaptcha

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"log"
	"math"
	"math/rand/v2"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

const (
	imageStringDpi = 72.0
	//MimeTypeImage output base64 mine-type.
	MimeTypeImage = "image/png"
)

// ItemChar captcha item of unicode characters
type ItemChar struct {
	bgColor color.Color
	width   int
	height  int
	nrgba   *image.NRGBA
}

// NewItemChar creates a captcha item of characters
func NewItemChar(w int, h int, bgColor color.RGBA) *ItemChar {
	d := ItemChar{width: w, height: h}
	m := image.NewNRGBA(image.Rect(0, 0, w, h))
	draw.Draw(m, m.Bounds(), &image.Uniform{C: bgColor}, image.Point{}, draw.Src)
	d.nrgba = m
	return &d
}

// drawHollowLine draw strong and bold white line.
func (item *ItemChar) drawHollowLine() *ItemChar {
	first := item.width / 20
	end := first * 19

	lineColor := RandLightColor()

	x1 := float64(rand.IntN(first))
	x2 := float64(rand.IntN(first) + end)

	multiple := float64(rand.IntN(5)+3) / float64(5)
	if int(multiple*10)%3 == 0 {
		multiple = multiple * -1.0
	}

	w := item.height / 20

	for ; x1 < x2; x1++ {
		y := math.Sin(x1*math.Pi*multiple/float64(item.width)) * float64(item.height/3)
		if multiple < 0 {
			y = y + float64(item.height/2)
		}
		item.nrgba.Set(int(x1), int(y), lineColor)
		for i := 0; i <= w; i++ {
			item.nrgba.Set(int(x1), int(y)+i, lineColor)
		}
	}

	return item
}

// drawSineLine draw a sine line.
func (item *ItemChar) drawSineLine() *ItemChar {
	var py float64

	a := rand.IntN(item.height / 2)
	b := random(int64(-item.height/4), int64(item.height/4))
	f := random(int64(-item.height/4), int64(item.height/4))

	var t float64
	if item.height > item.width/2 {
		t = random(int64(item.width/2), int64(item.height))
	} else if item.height == item.width/2 {
		t = float64(item.height)
	} else {
		t = random(int64(item.height), int64(item.width/2))
	}
	w := float64((2 * math.Pi) / t)

	px1 := 0
	px2 := int(random(int64(float64(item.width)*0.8), int64(item.width)))

	c := RandDeepColor()
	cNRGBA := color.NRGBA{R: c.R, G: c.G, B: c.B, A: c.A}

	for px := px1; px < px2; px++ {
		if w != 0 {
			py = float64(a)*math.Sin(w*float64(px)+f) + b + (float64(item.width) / float64(5))
			i := item.height / 5
			for i > 0 {
				item.nrgba.SetNRGBA(px+i, int(py), cNRGBA)
				i--
			}
		}
	}

	return item
}

// drawSlimLine draw n slim-random-color lines.
func (item *ItemChar) drawSlimLine(num int) *ItemChar {
	first := item.width / 10
	end := first * 9
	y := item.height / 3

	for i := 0; i < num; i++ {
		point1 := point{X: rand.IntN(first), Y: rand.IntN(y)}
		point2 := point{X: rand.IntN(first) + end, Y: rand.IntN(y)}

		if i%2 == 0 {
			point1.Y = rand.IntN(y) + y*2
			point2.Y = rand.IntN(y)
		} else {
			point1.Y = rand.IntN(y) + y*(i%2)
			point2.Y = rand.IntN(y) + y*2
		}

		item.drawBeeline(point1, point2, RandDeepColor())
	}
	return item
}

func (item *ItemChar) drawBeeline(point1 point, point2 point, lineColor color.RGBA) {
	dx := math.Abs(float64(point1.X - point2.X))
	dy := math.Abs(float64(point2.Y - point1.Y))
	sx, sy := 1, 1
	if point1.X >= point2.X {
		sx = -1
	}
	if point1.Y >= point2.Y {
		sy = -1
	}
	err := dx - dy

	// تبدیل یک‌بار به NRGBA برای SetNRGBA سریع‌تر
	c := color.NRGBA{R: lineColor.R, G: lineColor.G, B: lineColor.B, A: lineColor.A}

	for {
		item.nrgba.SetNRGBA(point1.X, point1.Y, c)
		item.nrgba.SetNRGBA(point1.X+1, point1.Y, c)
		item.nrgba.SetNRGBA(point1.X-1, point1.Y, c)
		item.nrgba.SetNRGBA(point1.X+2, point1.Y, c)
		item.nrgba.SetNRGBA(point1.X-2, point1.Y, c)
		if point1.X == point2.X && point1.Y == point2.Y {
			return
		}
		e2 := err * 2
		if e2 > -dy {
			err -= dy
			point1.X += sx
		}
		if e2 < dx {
			err += dx
			point1.Y += sy
		}
	}
}

func (item *ItemChar) drawNoise(noiseText string, fonts []*truetype.Font) error {
	c := freetype.NewContext()
	c.SetDPI(imageStringDpi)
	c.SetClip(item.nrgba.Bounds())
	c.SetDst(item.nrgba)
	c.SetHinting(font.HintingFull)

	// فونت را یک‌بار بیرون حلقه تنظیم کن (اگر چند فونت داری، اینجا را تغییر بده)
	c.SetFont(fonts[0])

	rawFontSize := float64(item.height) / (1 + float64(rand.IntN(7))/float64(10))

	for _, char := range noiseText {
		rw := rand.IntN(item.width)
		rh := rand.IntN(item.height)
		fontSize := rawFontSize/2 + float64(rand.IntN(5))
		c.SetSrc(image.NewUniform(RandColor()))
		c.SetFontSize(fontSize)
		pt := freetype.Pt(rw, rh)
		if _, err := c.DrawString(string(char), pt); err != nil {
			log.Println(err)
		}
	}
	return nil
}

// drawText draw captcha string to image.
func (item *ItemChar) drawText(text string, fonts []*truetype.Font) error {
	if len(text) == 0 {
		return errors.New("text must not be empty")
	}

	f := fonts[0]
	fontWidth := item.width / len(text)
	colorRGBA := RandDarkColor()
	colorNRGBA := color.NRGBA{R: colorRGBA.R, G: colorRGBA.G, B: colorRGBA.B, A: 255}

	for i, s := range text {
		fontSize := item.height * (rand.IntN(5) + 9) / 16
		x := fontWidth*i + fontWidth/fontSize
		y := item.height/2 + fontSize/2 - rand.IntN(item.height/16*3)

		// ← گلیف از کش
		glyph := globalGlyphCache.getGlyph(f, s, fontSize)

		// ← رسم گلیف روی بوم با رنگ تصادفی
		drawGlyphColored(item.nrgba, glyph, x, y, colorNRGBA)
	}
	return nil
}

// drawGlyphColored گلیف (image.Alpha) را با رنگ مشخص روی بوم می‌کشد
func drawGlyphColored(dst *image.NRGBA, glyph *image.Alpha, x, y int, c color.NRGBA) {
	gb := glyph.Bounds()
	db := dst.Rect

	for gy := 0; gy < gb.Dy(); gy++ {
		for gx := 0; gx < gb.Dx(); gx++ {
			a := glyph.Pix[gy*glyph.Stride+gx]
			if a == 0 {
				continue
			}
			px := x + gx
			py := y + gy
			if px < db.Min.X || py < db.Min.Y || px >= db.Max.X || py >= db.Max.Y {
				continue
			}
			idx := (py-db.Min.Y)*dst.Stride + (px-db.Min.X)*4
			dst.Pix[idx+0] = c.R
			dst.Pix[idx+1] = c.G
			dst.Pix[idx+2] = c.B
			dst.Pix[idx+3] = a
		}
	}
}

// BinaryEncoding encodes an image to PNG and returns a byte slice.
func (item *ItemChar) BinaryEncoding() []byte {
	var buf bytes.Buffer
	if err := png.Encode(&buf, item.nrgba); err != nil {
		panic(err.Error())
	}
	return buf.Bytes()
}

// WriteTo writes captcha character in png format into the given io.Writer.
func (item *ItemChar) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(item.BinaryEncoding())
	return int64(n), err
}

// EncodeB64string با BestSpeed
func (item *ItemChar) EncodeB64string() string {
	var buf bytes.Buffer
	buf.Grow(8192)

	buf.WriteString("data:")
	buf.WriteString(MimeTypeImage)
	buf.WriteString(";base64,")

	enc := base64.NewEncoder(base64.StdEncoding, &buf)
	pngEnc := &png.Encoder{CompressionLevel: png.BestSpeed}
	if err := pngEnc.Encode(enc, item.nrgba); err != nil {
		panic(err.Error())
	}
	if err := enc.Close(); err != nil {
		panic(err.Error())
	}

	return buf.String()
}

type point struct {
	X int
	Y int
}

// randFontFrom choose random font family.
func randFontFrom(fonts []*truetype.Font) *truetype.Font {
	fontCount := len(fonts)
	index := rand.IntN(fontCount)
	return fonts[index]
}
