package auth

type ValueComparator interface {
	Compare(plaintext string, hashed string) (bool, error)
	Hash(plaintext string) (string, error)
}
