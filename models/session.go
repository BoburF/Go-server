package models

import "time"

type Session struct {
	ExpiresAt time.Time `json:"expiresAt" bson:"expiresAt"`
	ID        string    `json:"_id,omitempty" bson:"_id,omitempty"`
	Token     string    `json:"token" bson:"token"`
	UserID    string    `json:"userId" bson:"userId"`
}
