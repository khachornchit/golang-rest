package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string             `bson:"email" json:"email"`
	Name      string             `bson:"name" json:"name"`
	Password  string             `bson:"password,omitempty" json:"password"`
	CreatedAt time.Time          `bson:"create_at,omitempty" json:"create_at"`
}
