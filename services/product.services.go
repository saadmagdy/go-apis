package services

import (
	"basic_api/models"
	"basic_api/utils"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductServices interface {
	GetAllProducts() ([]*models.Product, error)
	CreateProduct(data *models.ProductCreate, uid string) error
	GetProductById(id string) (*models.Product, error)
	GetProductByName(name string) (*models.Product, error)
	UpdateProduct(data *models.ProductUpdate, pid, uid string) (int64, error)
	DeleteProduct(pid, uid string) (int64, error)
}
type ProductServicesImpl struct {
	productColl *mongo.Collection
	userColl    *mongo.Collection
}

func NewProductServices(pColl, uColl *mongo.Collection) ProductServices {
	return &ProductServicesImpl{
		productColl: pColl,
		userColl:    uColl,
	}
}

func (s *ProductServicesImpl) CreateProduct(data *models.ProductCreate, uid string) error {
	ctx, cansel := utils.Ctx()
	defer cansel()
	var productToSave *models.Product
	createdByID, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return err
	}
	productExsist, _ := s.GetProductByName(data.ProductName)
	if productExsist != nil {
		return errors.New("product name alleardy in use")
	}
	productToSave = &models.Product{
		ID:          primitive.NewObjectID(),
		ProductName: data.ProductName,
		Price:       data.Price,
		CreatedBy:   createdByID,
	}
	result, err := s.productColl.InsertOne(ctx, productToSave)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: createdByID}}
	update := bson.D{{Key: "$push", Value: bson.M{"products": result.InsertedID.(primitive.ObjectID)}}}
	_, err = s.userColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
func (s *ProductServicesImpl) GetProductByName(name string) (*models.Product, error) {
	ctx, cansel := utils.Ctx()
	defer cansel()
	var product *models.Product
	filter := bson.D{{Key: "product_name", Value: name}}
	if err := s.productColl.FindOne(ctx, filter).Decode(&product); err != nil {
		return nil, err
	}
	return product, nil
}
func (s *ProductServicesImpl) GetProductById(id string) (*models.Product, error) {
	ctx, cansel := utils.Ctx()
	defer cansel()
	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var product *models.Product
	filter := bson.D{{Key: "_id", Value: pid}}
	err = s.productColl.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		return nil, err
	}
	return product, nil
}
func (s *ProductServicesImpl) UpdateProduct(data *models.ProductUpdate, pid, uid string) (int64, error) {
	ctx, cansel := utils.Ctx()
	defer cansel()
	var oldProduct *models.Product
	var updatedProduct *models.Product

	productid, err := primitive.ObjectIDFromHex(pid)
	if err != nil {
		return 0, err
	}
	userid, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return 0, err
	}
	filter := bson.D{{Key: "_id", Value: productid}}
	err = s.productColl.FindOne(ctx, filter).Decode(&oldProduct)
	if err != nil {
		return 0, err
	}
	if userid != oldProduct.CreatedBy {
		return 0, errors.New("you are not the creator of this product")
	}
	updatedProduct = &models.Product{
		ID:          oldProduct.ID,
		ProductName: data.ProductName,
		Price:       data.Price,
		CreatedBy:   oldProduct.CreatedBy,
	}
	if data.ProductName == "" {
		updatedProduct.ProductName = oldProduct.ProductName
	} else if data.Price < 0 {
		updatedProduct.Price = oldProduct.Price
	}
	update := bson.D{{Key: "$set", Value: updatedProduct}}
	res, err := s.productColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}
func (s *ProductServicesImpl) DeleteProduct(pid, uid string) (int64, error) {
	ctx, cansel := utils.Ctx()
	defer cansel()
	var deletproduct *models.Product
	productid, err := primitive.ObjectIDFromHex(pid)
	if err != nil {
		return 0, err
	}
	userid, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return 0, err
	}
	filter := bson.D{{Key: "_id", Value: productid}}
	err = s.productColl.FindOne(ctx, filter).Decode(&deletproduct)
	if err != nil {
		return 0, err
	}
	if userid != deletproduct.CreatedBy {
		return 0, errors.New("you are not the creator of this product")
	}
	res, err := s.productColl.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}
	fltr := bson.D{{Key: "_id", Value: userid}}
	update := bson.D{{Key: "$pull", Value: bson.M{"products": productid}}}
	_, err = s.userColl.UpdateOne(ctx, fltr, update)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}
func (s *ProductServicesImpl) GetAllProducts() ([]*models.Product, error) {
	ctx, cansel := utils.Ctx()
	defer cansel()
	var products []*models.Product
	cur, err := s.productColl.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	for cur.Next(ctx) {
		var product *models.Product
		if err := cur.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	defer cur.Close(ctx)
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return products, nil
}
