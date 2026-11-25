package iranCaptcha

import (
	"image"
	"sync"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// glyphCache کش گلیف‌های rasterize شده
type glyphCache struct {
	mu    sync.RWMutex
	fonts map[*truetype.Font]map[rune]map[int]*image.Alpha // font -> rune -> size -> glyph
}

var globalGlyphCache = &glyphCache{
	fonts: make(map[*truetype.Font]map[rune]map[int]*image.Alpha),
}

// getGlyph گلیف rasterize شده را برمی‌گرداند (یا می‌سازد و کش می‌کند)
func (gc *glyphCache) getGlyph(f *truetype.Font, r rune, size int) *image.Alpha {
	gc.mu.RLock()
	if byRune, ok := gc.fonts[f]; ok {
		if bySize, ok := byRune[r]; ok {
			if glyph, ok := bySize[size]; ok {
				gc.mu.RUnlock()
				return glyph
			}
		}
	}
	gc.mu.RUnlock()

	// ساخت گلیف جدید
	glyph := rasterizeGlyph(f, r, size)

	gc.mu.Lock()
	if gc.fonts[f] == nil {
		gc.fonts[f] = make(map[rune]map[int]*image.Alpha)
	}
	if gc.fonts[f][r] == nil {
		gc.fonts[f][r] = make(map[int]*image.Alpha)
	}
	gc.fonts[f][r][size] = glyph
	gc.mu.Unlock()

	return glyph
}

// rasterizeGlyph یک گلیف را به image.Alpha تبدیل می‌کند
func rasterizeGlyph(f *truetype.Font, r rune, size int) *image.Alpha {
	// اندازه‌ی گلیف
	bounds := f.Bounds(fixed.I(size))
	width := bounds.Max.X.Ceil() - bounds.Min.X.Floor()
	height := bounds.Max.Y.Ceil() - bounds.Min.Y.Floor()
	if width <= 0 || height <= 0 {
		return image.NewAlpha(image.Rect(0, 0, 1, 1))
	}

	// کمی حاشیه
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

	// رسم در مبدأ
	pt := freetype.Pt(1, height-1)
	c.DrawString(string(r), pt)

	return dst
}
