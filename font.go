package iranCaptcha

import (
	"embed"
)

// defaultEmbeddedFontsFSF Built-in font storage.

//go:embed fonts/*.ttf
var defaultEmbeddedFontsFS embed.FS

var DefaultEmbeddedFonts = NewEmbeddedFontsStorage(defaultEmbeddedFontsFS)
