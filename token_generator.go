package auth

type TokenGenerator interface {
	GenerateToken(payload interface{}, secret string, expiresIn uint64) (string, error)
}
