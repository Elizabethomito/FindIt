package models

// User represents a registered user.
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"` // omit from responses
}
