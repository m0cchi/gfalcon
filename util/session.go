package util

import (
	"math/rand"
	"time"
)

var letters []rune

var letters_size int

func init() {
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_-")
	letters_size = len(letters)
}

func GenerateSessionID(size int) string {
	gfrand := rand.New(rand.NewSource(2332141))
	gfrand.Seed(time.Now().UnixNano())
	sessionID := make([]rune, size)
	for i := range sessionID {
		sessionID[i] = letters[gfrand.Intn(letters_size)]
	}
	return string(sessionID)
}
