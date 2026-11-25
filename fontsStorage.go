package iranCaptcha

import (
	"embed"
	"sync"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

type EmbeddedFontsStorage struct {
	fs    embed.FS
	cache sync.Map
}

func (s *EmbeddedFontsStorage) LoadFontByName(name string) *truetype.Font {
	if v, ok := s.cache.Load(name); ok {
		return v.(*truetype.Font)
	}

	fontBytes, err := s.fs.ReadFile(name)
	if err != nil {
		panic(err)
	}
	trueTypeFont, err := freetype.ParseFont(fontBytes)
	if err != nil {
		panic(err)
	}

	s.cache.Store(name, trueTypeFont)
	return trueTypeFont
}

func (s *EmbeddedFontsStorage) LoadFontsByNames(assetFontNames []string) []*truetype.Font {
	result := make([]*truetype.Font, 0, len(assetFontNames))
	for _, name := range assetFontNames {
		result = append(result, s.LoadFontByName(name))
	}
	return result
}

func NewEmbeddedFontsStorage(fs embed.FS) *EmbeddedFontsStorage {
	return &EmbeddedFontsStorage{fs: fs}
}
