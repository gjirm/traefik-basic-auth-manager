package tbam

import (
	"crypto/rand"
	"math/big"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Generate new random password
func generatePassword(length int, symbols string) (string, error) {
	result := ""
	for {
		if len(result) >= length {
			return result, nil
		}
		num, err := rand.Int(rand.Reader, big.NewInt(int64(127)))
		if err != nil {
			return "", err
		}
		s := string(rune(num.Int64()))
		if strings.Contains(symbols, s) {
			result += s
		}
	}
}

// Generate random bcrypt hash
func GetRandomBcryptHash() (password, hash string, err error) {

	password, err = generatePassword(12, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	if err != nil {
		return "", "", err
	}

	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	return password, string(passwordBytes), nil
}
