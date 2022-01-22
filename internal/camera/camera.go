package camera

import (
	"encoding/base64"
	"fmt"
	"image"
	"io/ioutil"

	"gocv.io/x/gocv"
)

// Snap opens the camera, takes a picture and returns the image data
func Snap() (string, error) {
	deviceID := 0
	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		return "", fmt.Errorf("error opening video capture device: %v", deviceID)
	}
	defer webcam.Close()

	mat := gocv.NewMat()
	defer mat.Close()

	if ok := webcam.Read(&mat); !ok {
		return "", fmt.Errorf("cannot read device %v", deviceID)
	}
	if mat.Empty() {
		return "", fmt.Errorf("no image on device %v", deviceID)
	}

	maxWidth := mat.Cols()
	middle := maxWidth / 2
	height := mat.Rows()
	square := image.Rect(0, middle-(height/2), middle+(height/2), mat.Rows())
	croppedMat := mat.Region(square)
	mat = croppedMat.Clone()

	out, err := gocv.IMEncode(gocv.PNGFileExt, mat)
	if err != nil {
		return "", err
	}

	gocv.IMWrite("cam1.jpg", mat)

	return base64.StdEncoding.EncodeToString(out), nil
}

func OpenImage(filename string) (string, error) {
	// mat := gocv.IMRead(filename, gocv.IMReadColor)
	// if mat.Empty() {
	// 	return "", fmt.Errorf("error reading image from: %v", filename)
	// }

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}
