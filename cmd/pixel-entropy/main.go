package main

import (
	"github.com/heshanthenura/pixel-entropy/internal/camera"
	"github.com/heshanthenura/pixel-entropy/internal/http"
)

func main() {
	go camera.StartEntropy()
	http.StartServer()
}
