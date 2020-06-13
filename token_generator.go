package auth

type TokenGenerator interface {
	GenerateToken(payload interface{}, secret string, expiresIn int64) (string, error)
}
