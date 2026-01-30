package avro

type UserCreatedAvro struct {
	UserID    string `avro:"user_id"`
	FirstName string `avro:"first_name"`
	LastName  string `avro:"last_name"`
	Email     string `avro:"email"`
	Avatar    string `avro:"avatar"`
}
