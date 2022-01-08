package camera

import (
	"encoding/base64"
	"io/ioutil"
)

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

// 	out, err := gocv.IMEncode(gocv.PNGFileExt, mat)
// 	if err != nil {
// 		return "", err
// 	}

// 	return base64.StdEncoding.EncodeToString(out), nil
// }

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
