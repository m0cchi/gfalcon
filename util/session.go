package util

import (
	"math/rand"
	"time"
)

var letters []rune

var lettersSize int

func init() {
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_-")
	lettersSize = len(letters)
}

func GenerateSessionID(size int) string {
	gfrand := rand.New(rand.NewSource(2332141))
	gfrand.Seed(time.Now().UnixNano())
	sessionID := make([]rune, size)
	for i := range sessionID {
		sessionID[i] = letters[gfrand.Intn(lettersSize)]
	}
	return string(sessionID)
}
