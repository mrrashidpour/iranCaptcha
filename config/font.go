package config

import (
	"io/ioutil"
	"log"

	"github.com/golang/freetype/truetype"
)

type FileFontStorage struct {
	Dir string
}

func (f *FileFontStorage) LoadFontByName(name string) *truetype.Font {
	path := f.Dir + "/" + name
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("Failed to read font:", err)
		return nil
	}
	tf, err := truetype.Parse(data)
	if err != nil {
		log.Println("Failed to parse font:", err)
		return nil
	}
	return tf
}

func (f *FileFontStorage) LoadFontsByNames(names []string) []*truetype.Font {
	fonts := []*truetype.Font{}
	for _, name := range names {
		if tf := f.LoadFontByName(name); tf != nil {
			fonts = append(fonts, tf)
		}
	}
	return fonts
}
