package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hashedPwd, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return hashedPwd, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	matched, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return matched, nil
}