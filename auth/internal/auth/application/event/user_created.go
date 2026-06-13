package event

type UserCreated struct {
	UserID    string `json:"user_id" avro:"user_id"`
	Email     string `json:"email" avro:"email"`
	Username  string `json:"username" avro:"username"`
	FirstName string `json:"first_name" avro:"first_name"`
	LastName  string `json:"last_name" avro:"last_name"`
}
