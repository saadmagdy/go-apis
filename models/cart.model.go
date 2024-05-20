package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Cart struct {
	ID       primitive.ObjectID   `json:"_id" bson:"_id"`
	UserID   primitive.ObjectID   `json:"user_id" bson:"user_id"`
	Products []primitive.ObjectID `json:"product_ids" bon:"products_ids"`
}

type CartClient struct {
	UserID   primitive.ObjectID   `json:"user_id" bson:"user_id"`
	Products []primitive.ObjectID `json:"product_ids" bon:"products_ids"`
}
