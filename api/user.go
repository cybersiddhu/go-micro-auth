package api

type UserJSON struct {
	Id        int    `json:"id"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}
