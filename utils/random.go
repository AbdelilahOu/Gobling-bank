package utils

import (
	"math/rand"
)

// func init() {
// 	rand.Seed(time.Now().UnixNano())
// }

func RandomInt(min, max int64) int64 {
	return rand.Int63n(max-min+1) + min
}

func RandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	s := make([]rune, length)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomAmount() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
