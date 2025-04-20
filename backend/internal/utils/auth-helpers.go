package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
)

const randomLettersAndDigits = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func HashPassword(password, salt string) string {
	data := password + salt
	hash := sha256.Sum256([]byte(data))

	return hex.EncodeToString(hash[:])
}

func GenSalt() string {
	res := make([]byte, 10)
	for i := 0; i < 10; i++ {
		res[i] = randomLettersAndDigits[rand.Intn(len(randomLettersAndDigits))]
	}

	return string(res)
}

func CheckPassword(password, userPassword, userSalt string) bool {
	passwordCheck := sha256.Sum256([]byte(password + userSalt))

	return hex.EncodeToString(passwordCheck[:]) == userPassword
}
