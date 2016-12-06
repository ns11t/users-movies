package model

// Represents Movie entity for API
type Movie struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Year    int    `json:"year"`
	GenreId int    `json:"-"`
	Genre   string `json:"genre"`
}
