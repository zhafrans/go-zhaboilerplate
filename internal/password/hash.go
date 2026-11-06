package password

import "golang.org/x/crypto/bcrypt"

// Hash hashes a plaintext password using bcrypt.
func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Check compares a hashed password with a plaintext candidate password.
func Check(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

