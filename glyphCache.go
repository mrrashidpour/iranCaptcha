package iranCaptcha

import (
	"image"
	"log"
	"sync/atomic"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type glyphCache struct {
	fonts map[*truetype.Font]map[rune]map[int]*image.Alpha
	ready atomic.Bool
}

var globalGlyphCache = &glyphCache{
	fonts: make(map[*truetype.Font]map[rune]map[int]*image.Alpha),
}

func (gc *glyphCache) putGlyph(f *truetype.Font, r rune, size int) {
	if gc.fonts[f] == nil {
		gc.fonts[f] = make(map[rune]map[int]*image.Alpha)
	}
	if gc.fonts[f][r] == nil {
		gc.fonts[f][r] = make(map[int]*image.Alpha)
	}
	gc.fonts[f][r][size] = rasterizeGlyph(f, r, size)
}

func (gc *glyphCache) getGlyph(f *truetype.Font, r rune, size int) *image.Alpha {
	if !gc.ready.Load() {
		panic("glyph cache not ready")
	}
	if byRune, ok := gc.fonts[f]; ok {
		if bySize, ok := byRune[r]; ok {
			if glyph, ok := bySize[size]; ok {
				return glyph
			}
		}
	}
	log.Printf("⚠️ glyph not pre-warmed: rune=%c size=%d", r, size)
	return image.NewAlpha(image.Rect(0, 0, 1, 1))
}

func (gc *glyphCache) markReady() {
	gc.ready.Store(true)
}

func rasterizeGlyph(f *truetype.Font, r rune, size int) *image.Alpha {
	// متریک‌های فونت
	fu := fixed.I(size)
	bounds := f.Bounds(fu)

	// ابعاد کادر
	width := bounds.Max.X.Ceil() - bounds.Min.X.Floor()
	height := bounds.Max.Y.Ceil() - bounds.Min.Y.Floor()
	if width <= 0 || height <= 0 {
		return image.NewAlpha(image.Rect(0, 0, 1, 1))
	}

	const pad = 3
	w := width + pad*2
	h := height + pad*2

	dst := image.NewAlpha(image.Rect(0, 0, w, h))

	c := freetype.NewContext()
	c.SetDPI(imageStringDpi)
	c.SetFont(f)
	c.SetFontSize(float64(size))
	c.SetClip(dst.Bounds())
	c.SetDst(dst)
	c.SetSrc(image.Opaque)
	c.SetHinting(font.HintingFull)

	// x: کادر از Min.X شروع می‌شود، پس با pad - Min.X شروع می‌کنیم
	x := pad - bounds.Min.X.Floor()

	// y: freetype از baseline به بالا رسم می‌کند و y از پایین‌بوم است.
	// کادر از baseline.Min.Y (معمولاً منفی) تا baseline.Max.Y (معمولاً مثبت) است.
	// پس baseline را در h - pad - bounds.Max.Y.Ceil() قرار می‌دهیم تا بالای کادر در pad بیفتد.
	y := h - pad - bounds.Max.Y.Ceil()

	pt := freetype.Pt(x, y)
	if _, err := c.DrawString(string(r), pt); err != nil {
		log.Printf("⚠️ DrawString error rune=%c size=%d: %v", r, size, err)
	}

	// بررسی خالی نبودن
	for _, v := range dst.Pix {
		if v != 0 {
			return dst
		}
	}
	log.Printf("⚠️ rasterizeGlyph empty rune=%c size=%d bounds=%v pt=%v", r, size, dst.Bounds(), pt)
	return dst
}
