package Models

type SignInResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}
