package core

import "math/rand"

const (
	letterBytes = "abcdefghijklmnoprstuwxyzABCDEFGHIJKLMNOPRSTUWXYZ"
)

func RandomFilename(n int) string {
	if n < 1 {
		return ""
	}

	b := make([]byte, n)

	for i := range n {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}
