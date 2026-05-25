package dto

type LoginRequest struct {
	EmailOrUsername string
	Password        string
	UserAgent       string
	IPAddress       string
}
