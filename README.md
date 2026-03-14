# pixel-entropy

A hardware random number generator that uses **real camera pixel noise** as its entropy source.

It continuously captures raw frames from your webcam and derives a SHA-256 hash directly from the full pixel data of each frame. The hash is stored in memory and served over HTTP, so any client can request a fresh random hash at any time.

<img src="demo.gif">

## How it works

1. Camera opens **once** on startup and runs in a background goroutine.
2. Every raw frame's full pixel bytes are hashed with SHA-256, no diffing, maximum data.
3. The latest camera hash is stored behind a **read-write mutex** (`sync.RWMutex`).
4. HTTP requests read from that store instantly, no camera open/close per request.
5. `/hash` returns a **unique derived hash per request** by mixing the latest camera hash with an atomic counter.
6. Multiple concurrent requests are safe and return different values.

## Prerequisites

| Requirement | Notes |
|---|---|
| **Go** 1.24+ | [https://go.dev/dl](https://go.dev/dl) |
| **OpenCV** 4.x | Required by GoCV |
| **GoCV** v0.43.0 | Go bindings for OpenCV — [https://gocv.io](https://gocv.io) |
| **Webcam** | Any `/dev/video0` compatible camera |

### Install OpenCV (Linux)

```bash
sudo apt-get install libopencv-dev
```

Or follow the full GoCV install guide: https://github.com/hybridgroup/gocv#how-to-install

## Run

```bash
git clone https://github.com/heshanthenura/pixel-entropy
cd pixel-entropy
go run cmd/pixel-entropy/main.go
```

The server starts on port **8080**.


## API

### `GET /`

Health check.

**Response**
```
OK
```

### `GET /hash`

Returns a unique SHA-256 value per request, derived from live camera entropy.

**Response `200 OK`**
```json
{
  "hash": "a3f1c29d4e..."
}
```

**Response `503 Service Unavailable`** *(camera not ready yet)*
```json
hash not ready yet, camera is still warming up
```

**Example**
```bash
curl http://localhost:8080/hash
```

**Concurrent requests example**
```bash
seq 1 5 | xargs -P5 -I{} curl -s http://localhost:8080/hash
```