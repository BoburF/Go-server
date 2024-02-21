package models

import "time"

type Transactions struct {
	CreatedTime time.Time `json:"createdTime" bson:"createdTime"`
	ID          string    `json:"_id,omitempty" bson:"_id,omitempty"`
	CategoryId  string    `json:"catgeoryId" bson:"categoryId"`
	UserId      string    `json:"userId" bson:"userId"`
	Amount      int       `json:"amount" bson:"amount"`
}
