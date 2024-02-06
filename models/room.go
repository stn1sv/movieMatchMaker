package models

type Room struct {
	Id                      string `json:"roomId"`
	Participants            []User `json:"participants"`
	Movies                  []Movie
	CreatedAt               int64
	MoviesThatEveryoneLiked []Movie
}
