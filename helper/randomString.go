package helper

import "math/rand"

const letters = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// CreateRandomString ..
func CreateRandomString(length int) string {
	buf := make([]byte, length)
	for i := range buf {
		buf[i] = letters[rand.Intn(len(letters))]
	}
	return string(buf)
}
