package support

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

func Now() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func Random() string {
	b := make([]byte, 32)
	if _, e := rand.Read(b); e != nil {
		panic(e)
	}
	return hex.EncodeToString(b)
}
