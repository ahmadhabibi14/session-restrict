package helper

import (
	"math/rand"
	"strconv"
)

func RandString(length int) string {
	const letterBytes = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890`

	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}

func GenerateOTP(digits int) (otp string) {
	for range digits {
		otp += strconv.Itoa(rand.Int() % 10)
	}

	return
}
