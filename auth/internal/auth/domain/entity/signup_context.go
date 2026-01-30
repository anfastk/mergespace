package entity

type SignupContextID string

type SignupContext struct {
	ID           SignupContextID
	Email        string
	FirstName    string
	LastName     string
	Username     string
	PasswordHash string
	OTP          string
}
