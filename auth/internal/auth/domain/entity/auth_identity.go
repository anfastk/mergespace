package entity

type AuthProvider string

const (
	AuthProviderGoogle AuthProvider = "google"
	AuthProviderGithub AuthProvider = "github"
	AuthProviderLocal  AuthProvider = "local"
)

type AuthIdentity struct {
	ID             string
	UserID         string
	Provider       AuthProvider
	ProviderUserID string
}
