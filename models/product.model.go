package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	ProductName string             `json:"product_name" bson:"product_name"`
	Price       float64            `json:"price" bson:"price"`
	CreatedBy   primitive.ObjectID `json:"created_by" bson:"created_by"`
}

type ProductCreate struct {
	ProductName string             `json:"product_name" bson:"product_name" validate:"required"`
	Price       float64            `json:"price" bson:"price" validate:"required"`
	CreatedBy   primitive.ObjectID `json:"created_by" bson:"created_by" validate:"required"`
}

type ProductUpdate struct {
	ProductName string  `json:"product_name" bson:"product_name" validate:"required"`
	Price       float64 `json:"price" bson:"price" validate:"required"`
}
