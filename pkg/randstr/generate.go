package randstr

import (
	"crypto/rand"
	"errors"
)

func Generate(size int) (string, error) {
	if size <= 0 {
		return "", errors.New("invalid size")
	}

	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	var str string
	for _, v := range b {
		str += string(chars[int(v)%len(chars)])
	}

	return str, nil
}
