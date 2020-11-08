package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

func main() {

	app := fiber.New()

	app.Get("/:companyName", func(c *fiber.Ctx) error {
		companyName := c.Params("companyName") // Get company name
		bg := c.Query("bg")

		imageType := c.Query("type")
		var image logo

		if imageType == "png" {
			image = logoPNG{companyName: processCompanyName(companyName)}
		} else {
			image = logoJPEG{companyName: processCompanyName(companyName)}
		}

		logo, _ := generate(image, bg) // Generate logo

		return c.SendStream(bytes.NewReader(logo))
	})

	app.Listen(":3000")
}

type resolution struct {
	width  int
	height int
}

type logoPNG struct {
	companyName string
}

type logoJPEG struct {
	companyName string
}

type logo interface {
	generateLogo(constColor string) ([]byte, error)
}

func generate(l logo, constColor string) ([]byte, error) {
	return l.generateLogo(constColor)
}

func processCompanyName(companyName string) string {
	companyName = strings.Replace(companyName, "_", " ", 1)
	return companyName
}

func selectResolution(companyNameLength int) resolution {
	res := resolution{}
	switch {
	case companyNameLength < 5:
		res = resolution{width: 1024, height: 768}
	case companyNameLength <= 10:
		res = resolution{width: 1366, height: 768}
	case companyNameLength < 15:
		res = resolution{width: 1920, height: 1020}
	default:
		res = resolution{width: 2560, height: 1440}
	}
	return res
}

func selectCenterPosition(companyNameLength int) int {
	// TODO Create an algorithm to handle this.

	var number int
	switch {
	case companyNameLength < 3:
		number = 600
	case companyNameLength == 3:
		number = 340
	case companyNameLength == 4:
		number = 200
	case companyNameLength == 5:
		number = 140
	case companyNameLength == 6:
		number = 90
	case companyNameLength == 7:
		number = 70
	case companyNameLength == 8:
		number = 90
	case companyNameLength == 9:
		number = 110
	case companyNameLength == 10:
		number = 115
	case companyNameLength == 13:
		number = 100
	case companyNameLength == 14:
		number = 110
	case companyNameLength < 15:
		number = 70
	case companyNameLength < 20:
		number = 70
	default:
		number = 40
	}
	return number
}

func selectFontSize(res resolution) int {
	var fontSize int
	switch {
	case res.width == 1024:
		fontSize = 160
	case res.width == 1366:
		fontSize = 155
	case res.width == 1920:
		fontSize = 170
	default:
		fontSize = 170
	}
	return fontSize
}

func (logo logoJPEG) generateLogo(constColor string) ([]byte, error) {

	// setup the background
	bgResolution := selectResolution(len(logo.companyName))

	upperLeft := image.Point{0, 0}
	lowerRight := image.Point{bgResolution.width, bgResolution.height}

	myLogo := image.NewRGBA(image.Rectangle{upperLeft, lowerRight})

	// background color
	backgroundColor := image.NewUniform(getRandomColor(constColor))

	// draw the background
	draw.Draw(myLogo, myLogo.Bounds(), backgroundColor, image.ZP, draw.Src)

	// write the company name
	fontSize := float64(selectFontSize(bgResolution))
	fontSpacing := float64(.7)

	myFontContext := getRandomFontAndContext()
	myFontContext.SetClip(myLogo.Bounds())
	myFontContext.SetDst(myLogo)
	myFontContext.SetFontSize(fontSize) //font size in points

	// center the logo
	companyNameLength := len(logo.companyName)
	centerPointX := (bgResolution.width / companyNameLength) + (bgResolution.width / companyNameLength) - selectCenterPosition(len(logo.companyName))

	pt := freetype.Pt(centerPointX, bgResolution.height/2+70)

	for _, str := range logo.companyName {
		stringLower := strings.ToLower(string(str))
		if stringLower == "i" || stringLower == "ı" {
			stringLower = "I"
		}

		_, err := myFontContext.DrawString(stringLower, pt)
		if err != nil {
			fmt.Println(err)
			return nil, nil
		}

		if stringLower == "I" || stringLower == "." {
			pt.X += myFontContext.PointToFixed(fontSize * (fontSpacing / 2))
		} else if stringLower == "w" || stringLower == "m" {
			pt.X += myFontContext.PointToFixed(fontSize * (fontSpacing * 1.3))
		} else if stringLower == "g" || stringLower == "o" || stringLower == "q" {
			pt.X += myFontContext.PointToFixed(fontSize * (fontSpacing * 1.1))
		} else if stringLower == "s" || stringLower == "e" {
			pt.X += myFontContext.PointToFixed(fontSize * (fontSpacing / 1.1))
		} else if stringLower == "l" {
			pt.X += myFontContext.PointToFixed(fontSize * (fontSpacing / 1.2))
		} else {
			pt.X += myFontContext.PointToFixed(fontSize * fontSpacing)
		}
	}

	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, myLogo, nil)
	if err != nil {
		panic(err)
	}

	return buf.Bytes(), nil

}

func (logo logoPNG) generateLogo(constColor string) ([]byte, error) {

	bgResolution := selectResolution(len(logo.companyName))

	upperLeft := image.Point{0, 0}
	lowerRight := image.Point{bgResolution.width, bgResolution.height}

	myLogo := image.NewRGBA(image.Rectangle{upperLeft, lowerRight})

	// write the company name
	fontSize := float64(selectFontSize(bgResolution))
	fontSpacing := float64(.7)

	myFontContext := getRandomFontAndContext()
	myFontContext.SetClip(myLogo.Bounds())
	myFontContext.SetDst(myLogo)
	myFontContext.SetFontSize(fontSize) //font size in points

	// center the logo
	companyNameLength := len(logo.companyName)
	centerPointX := (bgResolution.width / companyNameLength) + (bgResolution.width / companyNameLength) - selectCenterPosition(len(logo.companyName))

	pt := freetype.Pt(centerPointX, bgResolution.height/2+70)

	for _, str := range logo.companyName {
		stringLower := strings.ToLower(string(str))
		if stringLower == "i" || stringLower == "ı" {
			stringLower = "I"
		}

		_, err := myFontContext.DrawString(stringLower, pt)
		if err != nil {
			fmt.Println(err)
			return nil, nil
		}

		if stringLower == "I" || stringLower == "." {
			pt.X += myFontContext.PointToFixed(fontSize * (fontSpacing / 2))
		} else if stringLower == "w" || stringLower == "m" {
			pt.X += myFontContext.PointToFixed(fontSize * (fontSpacing * 1.3))
		} else if stringLower == "g" || stringLower == "o" || stringLower == "q" {
			pt.X += myFontContext.PointToFixed(fontSize * (fontSpacing * 1.1))
		} else if stringLower == "s" || stringLower == "e" {
			pt.X += myFontContext.PointToFixed(fontSize * (fontSpacing / 1.1))
		} else if stringLower == "l" {
			pt.X += myFontContext.PointToFixed(fontSize * (fontSpacing / 1.2))
		} else {
			pt.X += myFontContext.PointToFixed(fontSize * fontSpacing)
		}
	}

	buf := new(bytes.Buffer)
	err := png.Encode(buf, myLogo)
	if err != nil {
		panic(err)
	}

	return buf.Bytes(), nil

}

func getRandomFontAndContext() *freetype.Context {

	fontFile := fmt.Sprintf("fonts/%s.ttf", getRandomFont())
	fontDPI := float64(72)
	fontContext := new(freetype.Context)
	utf8Font := new(truetype.Font)

	fontColor := color.RGBA{255, 255, 255, 255}

	// TODO Fix font file
	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		log.Fatal(err)
	}

	utf8Font, err = freetype.ParseFont(fontBytes)
	if err != nil {
		log.Fatal(err)
	}

	fontForeGroundColor := image.NewUniform(fontColor)

	fontContext = freetype.NewContext()
	fontContext.SetDPI(fontDPI)
	fontContext.SetFont(utf8Font)
	fontContext.SetSrc(fontForeGroundColor)

	return fontContext
}

func getRandomColor(constColor string) color.RGBA {

	rand.Seed(time.Now().UnixNano())

	max := 125
	min := 0

	red := uint8(rand.Intn(max-min) + min)
	green := uint8(rand.Intn(max-min) + min)
	blue := uint8(rand.Intn(max-min) + min)
	if constColor != "" {
		return getConstantColor(constColor)
	}

	return color.RGBA{red, green, blue, 0xff}
}

func getRandomFont() string {
	fonts := []string{
		"Womby-Regular",
	}

	randomIndex := rand.Int() % len(fonts)

	return fonts[randomIndex]
}

type constantColors struct {
	name  string
	color color.RGBA
}

func getConstantColor(colorName string) color.RGBA {
	var constColor = color.RGBA{55, 0, 179, 255} // Set default color to Primary Variant #3700B3

	colors := []constantColors{
		constantColors{name: "red", color: color.RGBA{100, 12, 20, 255}},
		constantColors{name: "indigo", color: color.RGBA{63, 81, 181, 255}},
		constantColors{name: "teal", color: color.RGBA{0, 150, 136, 255}},
		constantColors{name: "brown", color: color.RGBA{121, 85, 72, 255}},
		constantColors{name: "deep orange", color: color.RGBA{255, 87, 34, 255}},
		constantColors{name: "gray", color: color.RGBA{158, 158, 158, 255}},
		constantColors{name: "primary", color: color.RGBA{55, 0, 179, 255}},
	}

	for _, color := range colors {
		if color.name == colorName {
			constColor = color.color
		}
	}
	return constColor
}
