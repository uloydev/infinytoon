package schema

type User struct {
	Base
	Name     string `bson:"name"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
	// ApiKey   *string `bson:"apiKey,omitempty"`
}

func (u *User) CollectionName() string {
	return "users"
}
