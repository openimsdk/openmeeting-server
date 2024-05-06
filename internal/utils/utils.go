package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"time"
)

func GenerateUniqueKey() string {
	return Md5(strconv.Itoa(int(time.Now().UnixMilli())))
}

func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	cipher := h.Sum(nil)
	return hex.EncodeToString(cipher)
}
