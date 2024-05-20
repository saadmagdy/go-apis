package services

import (
	"basic_api/models"
	"basic_api/utils"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CartServices interface {
	AddProductToCart(uid, pid string) error
	RemoveProductFromCart(uid, pid string) error
	GetUserCart(uid string) (*models.Cart, error)
}

type CartServicesImpl struct {
	cartColl    *mongo.Collection
	productColl *mongo.Collection
}

func NewCartServices(crtColl, productColl *mongo.Collection) CartServices {
	return &CartServicesImpl{
		cartColl:    crtColl,
		productColl: productColl,
	}
}

func (s *CartServicesImpl) AddProductToCart(uid, pid string) error {
	ctx, cansel := utils.Ctx()
	defer cansel()

	var cart *models.Cart
	var product *models.Product
	var newCart *models.Cart

	userId, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return err
	}
	productId, err := primitive.ObjectIDFromHex(pid)
	if err != nil {
		return err
	}

	productFilter := bson.D{{Key: "_id", Value: productId}}
	err = s.productColl.FindOne(ctx, productFilter).Decode(&product)
	if err == mongo.ErrNoDocuments {
		return errors.New("product not found to add to cart")
	} else if err != nil {
		return err
	}
	newCart = &models.Cart{
		ID:       primitive.NewObjectID(),
		UserID:   userId,
		Products: []primitive.ObjectID{productId},
	}
	userFilter := bson.D{{Key: "user_id", Value: userId}}
	err = s.cartColl.FindOne(ctx, userFilter).Decode(&cart)
	if err == mongo.ErrNoDocuments {
		_, err = s.cartColl.InsertOne(ctx, newCart)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	update := bson.D{{Key: "$push", Value: bson.M{"product_ids": productId}}}
	_, err = s.cartColl.UpdateOne(ctx, userFilter, update)
	if err != nil {
		return err
	}

	return nil
}
func (s *CartServicesImpl) RemoveProductFromCart(uid, pid string) error {
	ctx, cansel := utils.Ctx()
	defer cansel()
	var cart *models.Cart

	userId, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return err
	}
	productId, err := primitive.ObjectIDFromHex(pid)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "user_id", Value: userId}}
	err = s.cartColl.FindOne(ctx, filter).Decode(&cart)
	if err == mongo.ErrNoDocuments {
		return errors.New("this user has no cart")
	} else if err != nil {
		return err
	}

	for _, v := range cart.Products {
		if v != productId {
			return errors.New("the product is not in the user cart")
		}
		update := bson.D{{Key: "$pull", Value: bson.M{"product_ids": productId}}}
		_, err = s.cartColl.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}

	}

	return nil
}
func (s *CartServicesImpl) GetUserCart(uid string) (*models.Cart, error) {
	ctx, cansel := utils.Ctx()
	defer cansel()
	var cart *models.Cart
	userId, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "user_id", Value: userId}}
	err = s.cartColl.FindOne(ctx, filter).Decode(&cart)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("this user has no cart")
	} else if err != nil {
		return nil, err
	}
	return cart, nil
}
