package redis

type SignupContexModel struct {
	ID           string `json:"id"`
	FirstName    string `json:"firstname"`
	LastName     string `json:"lastname"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	OTP          string `json:"otp"`
}
