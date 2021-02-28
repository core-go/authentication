package auth

type TokenVerifier interface {
	VerifyToken(tokenString string, secret string) (map[string]interface{}, int64, int64, error)
}
