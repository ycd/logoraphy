package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"flag"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

func main() {
	var image logo

	companyName := flag.String("name", "", "Company name, a string.")

	imageType := flag.String("type", "jpeg", "Output format, JPEG or PNG.")
	bg := flag.String("bg", "", "Background color.")
	size := flag.Int("size", 0, "Pixel size of the output.")

	flag.Parse()

	if *companyName == "" {
		fmt.Println("Company name cannot be empty.")
		os.Exit(1)
	}

	if *imageType == "png" {
		image = logoPNG{companyName: processCompanyName(*companyName)}
	} else {
		image = logoJPEG{companyName: processCompanyName(*companyName)}
	}

	_ = generate(image, *bg, *size)

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
	generateLogo(constColor string, size int) error
}

func generate(l logo, constColor string, size int) error {
	return l.generateLogo(constColor, size)
}

func processCompanyName(companyName string) string {
	companyName = strings.Replace(companyName, "_", " ", 1)
	return companyName
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

func (logo logoJPEG) generateLogo(constColor string, size int) error {

	// setup the background
	bgResolution := selectResolution(len(logo.companyName))

	upperLeft := image.Point{0, 0}
	lowerRight := image.Point{bgResolution.width, bgResolution.height}

	myLogo := image.NewNRGBA(image.Rectangle{upperLeft, lowerRight})

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
			return nil
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

	f, err := os.Create(fmt.Sprintf("%s.jpeg", logo.companyName))
	if err != nil {
		panic(err)
	}

	defer f.Close()
	if size != 0 {
		myLogo = imaging.Resize(myLogo, size, 0, imaging.Lanczos) // Resize image with the given size
	}

	err = jpeg.Encode(f, myLogo, nil)
	if err != nil {
		panic(err)
	}

	return nil

}

func (logo logoPNG) generateLogo(constColor string, size int) error {

	bgResolution := selectResolution(len(logo.companyName))

	upperLeft := image.Point{0, 0}
	lowerRight := image.Point{bgResolution.width, bgResolution.height}

	myLogo := image.NewNRGBA(image.Rectangle{upperLeft, lowerRight})

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
			return nil
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

	f, err := os.Create(fmt.Sprintf("%s.png", logo.companyName))
	if err != nil {
		panic(err)
	}

	defer f.Close()
	if size != 0 {
		myLogo = imaging.Resize(myLogo, size, 0, imaging.Lanczos) // Resize image with the given size
	}

	err = png.Encode(f, myLogo)
	if err != nil {
		panic(err)
	}

	return nil

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
		constantColors{name: "black", color: color.RGBA{0,0,0,7}},
	}

	for _, color := range colors {
		if color.name == colorName {
			constColor = color.color
		}
	}
	return constColor
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
