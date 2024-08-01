package utils

import "golang.org/x/crypto/bcrypt"

func MatchHashtoPassword(originalHash string, password string) bool {
	inp_pwd := []byte(password)
	og_pwd := []byte(originalHash)
	err := bcrypt.CompareHashAndPassword(og_pwd, inp_pwd)
	return err == nil
}

func SaltNhash(pwd string) string {
	inp_pwd := []byte(pwd)
	hashedPwdBytes, err := bcrypt.GenerateFromPassword(inp_pwd, bcrypt.MinCost)
	CheckNilErr(err, "Unable to Hash password")
	return string(hashedPwdBytes)
}
