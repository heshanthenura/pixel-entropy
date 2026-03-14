package http

import (
	"encoding/json"
	nethttp "net/http"

	"github.com/heshanthenura/pixel-entropy/internal/hash"
)

func StartServer() {
	nethttp.HandleFunc("/", handler)
	nethttp.HandleFunc("/hash", hashHandler)
	nethttp.ListenAndServe(":8080", nil)
}

func handler(w nethttp.ResponseWriter, r *nethttp.Request) {
	w.Write([]byte("OK"))
}

func hashHandler(w nethttp.ResponseWriter, r *nethttp.Request) {
	if r.Method != nethttp.MethodGet {
		nethttp.Error(w, "method not allowed", nethttp.StatusMethodNotAllowed)
		return
	}

	h, err := hash.GetUniqueHash()
	if err != nil {
		nethttp.Error(w, err.Error(), nethttp.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"hash": h})
}
