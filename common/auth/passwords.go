package auth

import "golang.org/x/crypto/bcrypt"

const PWDSALT = "0dTeybIaJskijshaglve"

func HashPwd(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password+PWDSALT), bcrypt.DefaultCost+1)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
