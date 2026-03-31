package models

type Item struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Type        string `json:"type"` // lost or found
	Name        string `json:"name"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Date        string `json:"date"`
}
