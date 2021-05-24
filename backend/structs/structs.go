package structs

type ShoppingList struct {
	Title        string   `json:"title"`
	Items        []string `json:"items"`
	Participants string   `json:"participants"`
	Owner        *User    `json:"owner"`
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
