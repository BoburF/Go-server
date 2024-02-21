package models

type User struct {
	ID      string `json:"_id,omitempty" bson:"_id,omitempty"`
	Name    string `json:"name" bson:"name"`
	Surname string `json:"surname" bson:"surname"`
}
