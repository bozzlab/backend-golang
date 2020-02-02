package main

import (
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"

	"github.com/disintegration/imaging"
)

func main() {
	http.HandleFunc("/", resizeHandler)
	http.ListenAndServe(":3333", nil)
}

func resizeHandler(w http.ResponseWriter, r *http.Request) {
	// DO NOT accept other method than POST
	if r.Method != "POST" {
		http.Error(w, "Invalid Method", 400)
		return
	}
	if r.ContentLength <= 0 {
		http.Error(w, "require Image Data", 400)
		return
	}
	if r.ContentLength > 2*1024*1024 {
		http.Error(w, "Invalid Lenght", 400)
		return
	}
	// limit body to 2 MiB
	bodyImage := io.LimitReader(r.Body, 2*1024*1024)

	// r.Body.Read()

	// decode image using image.Decode
	img, imgType, err := image.Decode(bodyImage)
	if err != nil {
		http.Error(w, "Invalid Image", 400)

		return
	}
	// print image type (jpg, png, bmp) to console
	log.Print(imgType)

	// use imaging.Thumbnail to resize image
	img = imaging.Thumbnail(img, 150, 150, imaging.Gaussian)

	// encode result image using jpeg.Encode
	jpeg.Encode(w, img, &jpeg.Options{Quality: 80})
	// write result jpeg image to responseWriter
}
