package models

import (
	"math/rand"
	"time"
)

/*
	This function generates random combinations of characters.
	The allCases and lowerCases constants contains all the valid characters
*/

//Character sets used for string generation.
//lowerCases is used when requiresLowerCases is true
const allCases = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
const lowerCases = "abcdefghijklmnopqrstuvwxyz123456789"

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

//RandStringBytesMaskImprSrc generates and returns random indexes.
func RandStringBytesMaskImprSrc(n int, requiresLowerCases bool) string {
	var letterBytes string
	if requiresLowerCases {
		letterBytes = lowerCases
	} else {
		letterBytes = allCases
	}

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
