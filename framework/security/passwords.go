package security

import (
	"crypto/hmac"
	"encoding/hex"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
)

func PasswordOK(hash, password string) bool {
	if strings.HasPrefix(hash, "$2") {
		return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
	}
	parts := strings.Split(hash, "$")
	if len(parts) != 3 {
		return false
	}
	params := strings.Split(parts[0], ":")
	if len(params) != 4 || params[0] != "scrypt" {
		return false
	}
	n, _ := strconv.Atoi(params[1])
	r, _ := strconv.Atoi(params[2])
	p, _ := strconv.Atoi(params[3])
	if n < 2 || n > 1048576 || r < 1 || r > 32 || p < 1 || p > 16 {
		return false
	}
	b, e := scrypt.Key([]byte(password), []byte(parts[1]), n, r, p, 64)
	return e == nil && hmac.Equal([]byte(hex.EncodeToString(b)), []byte(parts[2]))
}
