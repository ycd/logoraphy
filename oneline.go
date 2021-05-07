package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"log"
	"strings"

	"github.com/fogleman/gg"
)

// OneLetter creates a new a new typographic favicon.
func OneLetter(letter string) (image *bytes.Buffer) {
	const S = 256

	dc := gg.NewContext(S, S)
	dc.SetRGB(0.4039, 0.34901, 1)
	dc.Clear()
	dc.SetRGB(1, 1, 1)
	if err := dc.LoadFontFace("font.ttf", 201); err != nil {
		panic(err)
	}

	dc.DrawRoundedRectangle(0, 0, 216, 512, 0)

	dc.DrawStringAnchored(letter, S/2, S/2, 0.5, 0.5)
	dc.Clip()
	buf := new(bytes.Buffer)
	err := dc.EncodePNG(buf)
	if err != nil {
		log.Fatal(err)
	}

	return buf
}


