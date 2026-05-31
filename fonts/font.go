package fonts

import (
	"embed"

	"github.com/mrrashidpour/iranCaptcha/config"
)

// defaultEmbeddedFontsFSF Built-in font storage.

//go:embed *.ttf
var defaultEmbeddedFontsFS embed.FS

var DefaultEmbeddedFonts = config.NewEmbeddedFontsStorage(defaultEmbeddedFontsFS)
