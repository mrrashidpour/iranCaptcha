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

// glyphCache کش گلیف‌های rasterize شده (بدون قفل بعد از pre-warm)
type glyphCache struct {
	fonts map[*truetype.Font]map[rune]map[int]*image.Alpha
	ready atomic.Bool
}

var globalGlyphCache = &glyphCache{
	fonts: make(map[*truetype.Font]map[rune]map[int]*image.Alpha),
}

// putGlyph گلیف را می‌سازد و در کش می‌گذارد (فقط در pre-warm)
func (gc *glyphCache) putGlyph(f *truetype.Font, r rune, size int) {
	if gc.fonts[f] == nil {
		gc.fonts[f] = make(map[rune]map[int]*image.Alpha)
	}
	if gc.fonts[f][r] == nil {
		gc.fonts[f][r] = make(map[int]*image.Alpha)
	}
	gc.fonts[f][r][size] = rasterizeGlyph(f, r, size)
}

// getGlyph گلیف rasterize شده را برمی‌گرداند (بدون قفل — بعد از pre-warm)
func (gc *glyphCache) getGlyph(f *truetype.Font, r rune, size int) *image.Alpha {
	if !gc.ready.Load() {
		panic("glyph cache not ready (pre-warm not called)")
	}

	if byRune, ok := gc.fonts[f]; ok {
		if bySize, ok := byRune[r]; ok {
			if glyph, ok := bySize[size]; ok {
				return glyph
			}
		}
	}

	// اگر به اینجا رسید، pre-warm ناقص بوده
	log.Printf("⚠️ glyph not pre-warmed: rune=%c size=%d — falling back to empty", r, size)
	return image.NewAlpha(image.Rect(0, 0, 1, 1))
}

// markReady بعد از اتمام pre-warm صدا زده می‌شود
func (gc *glyphCache) markReady() {
	gc.ready.Store(true)
}

// rasterizeGlyph یک گلیف را به image.Alpha تبدیل می‌کند
func rasterizeGlyph(f *truetype.Font, r rune, size int) *image.Alpha {
	bounds := f.Bounds(fixed.I(size))
	width := bounds.Max.X.Ceil() - bounds.Min.X.Floor()
	height := bounds.Max.Y.Ceil() - bounds.Min.Y.Floor()
	if width <= 0 || height <= 0 {
		return image.NewAlpha(image.Rect(0, 0, 1, 1))
	}

	width += 2
	height += 2

	dst := image.NewAlpha(image.Rect(0, 0, width, height))

	c := freetype.NewContext()
	c.SetDPI(imageStringDpi)
	c.SetFont(f)
	c.SetFontSize(float64(size))
	c.SetClip(dst.Bounds())
	c.SetDst(dst)
	c.SetSrc(image.Opaque)
	c.SetHinting(font.HintingFull)

	pt := freetype.Pt(1, height-1)
	c.DrawString(string(r), pt)

	return dst
}
