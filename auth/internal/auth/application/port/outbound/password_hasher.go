package outbound

type PasswordHasher interface {
	Hash(plain string) (string, error)
	Compare(hash string, plain string) error
}
