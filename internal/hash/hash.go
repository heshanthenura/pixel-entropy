package hash

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"sync"
	"sync/atomic"
)

var (
	mu      sync.RWMutex
	current [32]byte
	ready   bool
	seq     atomic.Uint64
)

func StoreHash(h [32]byte) {
	mu.Lock()
	defer mu.Unlock()
	current = h
	ready = true
}

func GetHash() (string, error) {
	mu.RLock()
	defer mu.RUnlock()
	if !ready {
		return "", errors.New("hash not ready yet, camera is still warming up")
	}
	return hex.EncodeToString(current[:]), nil
}

func GetUniqueHash() (string, error) {
	mu.RLock()
	if !ready {
		mu.RUnlock()
		return "", errors.New("hash not ready yet, camera is still warming up")
	}
	base := current
	mu.RUnlock()

	nonce := seq.Add(1)
	buf := make([]byte, 40)
	copy(buf[:32], base[:])
	binary.LittleEndian.PutUint64(buf[32:], nonce)

	derived := sha256.Sum256(buf)
	return hex.EncodeToString(derived[:]), nil
}
