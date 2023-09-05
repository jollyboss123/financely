package authentication

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" validation:"required"`
	Password  string `json:"password" validation:"required,gt=13"`
}

type LoginRequest struct {
	Email    string `json:"email" validation:"required"`
	Password string `json:"password" validation:"required"`
}
