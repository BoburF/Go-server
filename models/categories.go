package models

type Categories struct {
	ID    string `json:"_id,omitempty" bson:"_id,omitempty"`
	Title string `json:"title" bson:"title"`
}
