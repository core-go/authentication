package auth

type PasswordDecrypter interface {
	Decrypt(cipherText string, secretKey string) (string, error)
}
