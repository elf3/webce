package password

import "golang.org/x/crypto/bcrypt"

// 加密密码
func Encrypt(source string) (string, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(source), bcrypt.DefaultCost)
	return string(password), err
}

// 对比密码
func Compare(hashPass, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(password))
}
