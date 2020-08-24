package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Person struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Nombre   string             `bson:"name" json:"nombre,omitempty"`
	Insultos int                `bson:"curses" json:"insultos,omitempty"`
}
