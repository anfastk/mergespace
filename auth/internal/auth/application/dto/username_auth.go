package dto

type CheckUsernameReq struct {
	Username string
}

type CheckUsernameRes struct {
	Available   bool
	Message     string
	Suggestions []string
}