package auth

type CodeGenerator interface {
	Generate() string
}
