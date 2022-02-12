package camera

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
)

var binary string

func init() {
	var err error
	binary, err = exec.LookPath("fswebcam")
	if err != nil {
		panic(fmt.Errorf("fswebcam not installed: %v", err))
	}
}

// Snap opens the camera, takes a picture and returns the image data
// func Snap() (string, error) {
// 	deviceID := 0
// 	webcam, err := gocv.OpenVideoCapture(deviceID)
// 	if err != nil {
// 		return "", fmt.Errorf("error opening video capture device: %v", deviceID)
// 	}
// 	defer webcam.Close()

// 	mat := gocv.NewMat()
// 	defer mat.Close()

// 	if ok := webcam.Read(&mat); !ok {
// 		return "", fmt.Errorf("cannot read device %v", deviceID)
// 	}
// 	if mat.Empty() {
// 		return "", fmt.Errorf("no image on device %v", deviceID)
// 	}

// 	maxWidth := mat.Cols()
// 	middle := maxWidth / 2
// 	height := mat.Rows()
// 	square := image.Rect(0, middle-(height/2), middle+(height/2), mat.Rows())
// 	croppedMat := mat.Region(square)
// 	defer croppedMat.Close()
// 	squareMat := croppedMat.Clone()
// 	defer squareMat.Close()

// 	out, err := gocv.IMEncode(gocv.PNGFileExt, squareMat)
// 	if err != nil {
// 		return "", err
// 	}

// 	gocv.IMWrite("cam1.jpg", squareMat)

// 	return base64.StdEncoding.EncodeToString(out.GetBytes()), nil
// }

func Snap() (string, error) {
	filename := "snap.jpg"

	args := []string{"fswebcam", "--delay", "2", "--skip", "20", "--no-banner", filename}
	env := os.Environ()

	err := syscall.Exec(binary, args, env)
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
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
