package says

import (
	"bytes"
	"image/color"
	"image/png"
	"io"
	"log"

	"github.com/fogleman/gg"
)

func VVSays(quote string) (io.Reader, error) {
	baseImg, err := gg.LoadImage("./vvimg.jpg")
	if err != nil {
		log.Printf("failed to load image: %v", err)
		return nil, err
	}

	imgWidth := baseImg.Bounds().Dx()
	imgHeight := baseImg.Bounds().Dy()

	dc := gg.NewContext(imgWidth, imgHeight)
	dc.DrawImage(baseImg, 0, 0)

	if err := dc.LoadFontFace("COMIC.TTF", 32); err != nil {
		log.Printf("error loading font face: %v", err)
		return nil, err
	}

	x := float64(178)
	y := float64(465)
	maxWidth := float64(380)

	dc.SetColor(color.Black)
	text := quote
	dc.DrawStringWrapped(text, x, y, 0, 0, maxWidth, 1.5, gg.AlignCenter)

	imgafter := dc.Image()
	buf := new(bytes.Buffer)
	err = png.Encode(buf, imgafter)
	if err != nil {
		return nil, err
	}

	data := buf.Bytes()
	reader := bytes.NewReader(data)

	return reader, nil

}
