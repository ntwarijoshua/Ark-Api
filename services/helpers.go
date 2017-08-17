package services

import (
	"math/rand"
	"time"
	"strconv"
	"golang.org/x/crypto/bcrypt"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func GenerateApiKey() string {
	len := 10
	characters := "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return StringWithCharset(len, characters)
}

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func ConvertParametersToIntegers(param string) int {
	x, _ := strconv.ParseInt(param, 0, 8)
	y := int(x)
	return y
}

func HashPassword(password string) string {
	pswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if (err != nil) {
		panic("Hashing Password Failed")
	}
	stringPwd := string(pswd)
	return stringPwd
}
