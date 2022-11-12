package oauth2

type TokenPort interface {
	GenerateToken(payload interface{}, secret string, expiresIn int64) (string, error)
	VerifyToken(tokenString string, secret string) (map[string]interface{}, int64, int64, error)
}
