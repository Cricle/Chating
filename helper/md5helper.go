package helper

import (
	"crypto/md5"
	"encoding/hex"
)

func MakeMd5Str(val string) (string, error) {
	m5 := md5.New()
	_, err := m5.Write([]byte(val))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(m5.Sum(nil)), nil
}
