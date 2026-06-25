package dto

type UpdateProfileRequest struct {
	UserID    string  `json:"-"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Bio       *string `json:"bio"`
	AvatarURL *string `json:"avatar_url"`
}
