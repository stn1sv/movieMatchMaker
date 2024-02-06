package models

type Input struct {
	Genres []string `json:"genres"`
	Years  []string `json:"years"`
}

type Movie struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	EnName      string `json:"enName"`
	Year        int    `json:"year"`
	Description string `json:"description"`
	Rating      struct {
		Kp   float64 `json:"kp"`
		Imdb float64 `json:"imdb"`
	} `json:"rating"`
	AgeRating int `json:"ageRating"`
	Poster    struct {
		URL string `json:"url"`
	} `json:"poster"`
	Genres []struct {
		Name string `json:"name"`
	} `json:"genres"`
	Countries []struct {
		Name string `json:"name"`
	} `json:"countries"`
	UsersWhoLikedItCount int `json:"-"`
}
