package auxiliary

import (
	"bytes"
	"image/jpeg"
	"log"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

func AphrConvertJpeg2WebP(jpegBytes []byte, webpBytes chan []byte) {
	jpegReader := bytes.NewReader(jpegBytes)
	jpegImage, errJpeg := jpeg.Decode(jpegReader)
	if errJpeg != nil {
		log.Fatalln("Error jpeg decode: ", errJpeg)
	}

	resizedImage := resize.Resize(512, 512, jpegImage, resize.Lanczos3)

	webpBuffer := bytes.NewBuffer(nil)
	errWebp := webp.Encode(webpBuffer, resizedImage, &webp.Options{Quality: 100, Lossless: true})
	if errWebp != nil {
		log.Fatalln("Error encode webp: ", errWebp)
	}

	webpBytes <- webpBuffer.Bytes()
}

func AphrConvertWebp2Jpeg(webpBytes []byte, jpegBytes chan []byte) {

	webpReader := bytes.NewReader(webpBytes)
	webpImg, errWebp := webp.Decode(webpReader)
	if errWebp != nil {
		log.Fatalln("Error convert webp: ", errWebp)
	}

	jpegBuffer := bytes.NewBuffer(nil)
	errJpeg := jpeg.Encode(jpegBuffer, webpImg, &jpeg.Options{Quality: 100})

	if errJpeg != nil {
		log.Fatalln("Error convert jpeg: ", errJpeg)
	}

	jpegBytes <- jpegBuffer.Bytes()

}

func AphrConvertMp4toWebp(mp4Bytes []byte, webpBytes chan []byte) {
	mp4Reader := bytes.NewReader(mp4Bytes)
	ff := ffmpeg_go.Input("pipe:").WithInput(mp4Reader)

	webpProcessor := bytes.NewBuffer(nil)
	webpCv := ff.Output("pipe:", ffmpeg_go.KwArgs{
		"ss":                "1",
		"t":                 "3",
		"f":                 "webp",
		"filter:v":          "fps=fps=10",
		"lossless":          "0",
		"loop":              "0",
		"vsync":             "0",
		"compression_level": "6",
		"q:v":               "10",
		"s":                 "512:512",
	})

	if errWebp := webpCv.WithOutput(webpProcessor).Run(); errWebp != nil {
		log.Fatalln("Error convert to sticker: ", errWebp)
	}

	webpBytes <- webpProcessor.Bytes()
}

func TestResizeImage(jpegBytes []byte, resultBytes []byte) {
	jpegReader := bytes.NewReader(jpegBytes)

	jpgImg, errJpg := jpeg.Decode(jpegReader)
	if errJpg != nil {
		log.Fatalln("Error decoding jpeg: ", errJpg)
	}

	widthImg := jpgImg.Bounds().Dx()
	heightImg := jpgImg.Bounds().Dy()

	options := &jpeg.Options{Quality: 3}
	resize.Resize(uint(widthImg), uint(heightImg), jpgImg, resize.InterpolationFunction(options.Quality))

}
