package util

import (
	"math/rand"
	"time"
)

var letters []rune

var letters_size int

var gfrand *rand.Rand

func init() {
	gfrand = rand.New(rand.NewSource(2332141))
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_-")
	letters_size = len(letters)
	gfrand = rand.New(rand.NewSource(1))
	gfrand.Seed(time.Now().UnixNano())
}

func GenerateSessionID(size int) string {
	sessionID := make([]rune, size)
	for i := range sessionID {
		sessionID[i] = letters[gfrand.Intn(letters_size)]
	}
	return string(sessionID)
}
