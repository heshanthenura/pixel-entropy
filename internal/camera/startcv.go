package camera

import (
	"crypto/sha256"
	"fmt"

	hashstore "github.com/heshanthenura/pixel-entropy/internal/hash"
	"gocv.io/x/gocv"
)

func openCamera(number int) (*gocv.VideoCapture, error) {
	webcam, err := gocv.OpenVideoCapture(number)
	if err != nil {
		return nil, fmt.Errorf("cannot open camera: %v", err)
	}
	if !webcam.IsOpened() {
		webcam.Close()
		return nil, fmt.Errorf("camera %d is not available", number)
	}
	return webcam, nil
}

func readFrame(webcam *gocv.VideoCapture, frame *gocv.Mat) bool {
	if ok := webcam.Read(frame); !ok || frame.Empty() {
		fmt.Println("Error reading frame from camera")
		return false
	}
	return true
}

func hashFrame(frame *gocv.Mat, debug bool) ([32]byte, bool) {
	var hash [32]byte

	data, err := frame.DataPtrUint8()
	if err != nil {
		fmt.Println("Error reading frame data:", err)
		return hash, false
	}

	hash = sha256.Sum256(data)
	if debug {
		fmt.Printf("Hash: %x\n", hash)
	}
	return hash, true
}

func generateHash(webcam *gocv.VideoCapture, frame *gocv.Mat, window *gocv.Window, showWindow bool, debug bool) string {
	var hash [32]byte
	for {
		if !readFrame(webcam, frame) {
			continue
		}

		if showWindow {
			window.IMShow(*frame)
		}

		currentHash, ok := hashFrame(frame, debug)
		if !ok {
			continue
		}
		hash = currentHash
		hashstore.StoreHash(currentHash)

		if window.WaitKey(1) == 27 {
			break
		}
	}
	return fmt.Sprintf("%x", hash)
}

func StartEntropy() {

	webcam, err := openCamera(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer webcam.Close()

	window := gocv.NewWindow("Frame Difference")
	defer window.Close()

	frame := gocv.NewMat()
	defer frame.Close()

	hash := generateHash(webcam, &frame, window, true, true)
	if hash != "" {
		fmt.Printf("%s\n", hash)
	}
}
