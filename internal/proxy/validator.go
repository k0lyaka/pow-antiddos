package proxy

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"time"
)

const MIN_NONCE_SIZE = 8
const MAX_NONCE_SIZE = 32

type ValidationRequest struct {
	Prefix     string
	Nonce      string
	Difficulty int
}

func checkTimestamp(nonce []byte) bool {
	timestamp := binary.BigEndian.Uint32(nonce[:4])

	return time.Now().UTC().Unix()-int64(timestamp) < 60
}

func checkHash(hash []byte, difficulty int) bool {
	offset := 0

	for i := 0; i <= difficulty-8; i += 8 {
		if hash[offset] != 0 {
			return false
		}
		offset++
	}

	mask := byte(0xff << (8 + difficulty - (offset * 8)))

	return hash[offset]&mask == 0
}

func Validate(req ValidationRequest) bool {
	if len(req.Nonce) < MIN_NONCE_SIZE || len(req.Nonce) > MAX_NONCE_SIZE {
		return false
	}

	// convert nonce to bytes
	nonce, err := hex.DecodeString(req.Nonce)
	if err != nil {
		return false
	}

	// create hash
	h := sha256.New()
	h.Write([]byte(req.Prefix))
	h.Write(nonce)
	hashed := h.Sum(nil)

	return checkTimestamp(nonce) && checkHash(hashed, req.Difficulty)
}
