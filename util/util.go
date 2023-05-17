package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandInt generates a integer between min and max
func RandInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandString generates a random string of length n
func RandString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandOwner generates a random owner name
func RandOwner() string {
	return RandString(8)
}

// RandBalance generates a random amount of money
func RandBalance() int64 {
	return RandInt(0, 1000)
}

// RandCurrency generates a random kind of currency
func RandCurrency() string {
	currencies := []string{EUR, CAD, USD}
	return currencies[rand.Intn(len(currencies))]
}
