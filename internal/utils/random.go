package utils

import "golang.org/x/exp/rand"

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
			b[i] = letterBytes[rand.Int63() % int64(len(letterBytes))]
	}
	return string(b)
}
