package securetools

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
)

func generateSalt() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)

	salt := make([]byte, 16)
	for i := range salt {
		salt[i] = charset[rnd.Intn(len(charset))]
	}
	return string(salt)
}

func HashPassword(password string) (string, string) {
	salt := generateSalt()
	hashed := md5.Sum([]byte(password + salt))
	return hex.EncodeToString(hashed[:]), salt
}

func VerifyPassword(password, salt string) string {
	hashed := md5.Sum([]byte(password + salt))
	return hex.EncodeToString(hashed[:])
}
