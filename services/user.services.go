package services

import (
	"basic_api/models"
	"basic_api/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServices interface {
	GetAllUsers() ([]*models.User, error)
	GetUserById(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.UserSignUp) error
	UpdateUser(id string, data *models.UserUpdate) (int64, error)
	DeleteUser(id string) (int64, error)
}

type UserServicesImpl struct {
	userColl *mongo.Collection
}

func NewUserService(userColl *mongo.Collection) UserServices {
	return &UserServicesImpl{
		userColl: userColl,
	}
}

func (s *UserServicesImpl) GetAllUsers() ([]*models.User, error) {
	ctx, cansel := utils.Ctx()
	defer cansel()
	filter := bson.D{{}}
	cur, err := s.userColl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var users []*models.User
	for cur.Next(ctx) {
		var user *models.User
		if err := cur.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = cur.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserServicesImpl) GetUserById(id string) (*models.User, error) {

	ctx, cansel := utils.Ctx()
	defer cansel()
	var user *models.User
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: objId}}
	err = s.userColl.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (s *UserServicesImpl) GetUserByEmail(email string) (*models.User, error) {
	ctx, cansel := utils.Ctx()
	defer cansel()
	var user *models.User
	filter := bson.D{{Key: "email", Value: email}}
	err := s.userColl.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (s *UserServicesImpl) CreateUser(user *models.UserSignUp) error {
	ctx, cansel := utils.Ctx()
	defer cansel()
	pass, err := utils.HashPWD(user.Password)
	if err != nil {
		return err
	}
	var newUser models.User
	newUser.ID = primitive.NewObjectID()
	newUser.Email = user.Email
	newUser.Name = user.Name
	newUser.Password = pass
	newUser.User_Type = user.User_Type
	_, err = s.userColl.InsertOne(ctx, newUser)
	if err != nil {
		return err
	}
	return nil
}
func (s *UserServicesImpl) UpdateUser(id string, data *models.UserUpdate) (int64, error) {
	ctx, cansel := utils.Ctx()
	defer cansel()
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}
	filter := bson.D{{Key: "_id", Value: objId}}

	var foundUser *models.User
	var userToUpdate models.User

	if err := s.userColl.FindOne(ctx, filter).Decode(&foundUser); err != nil {
		return 0, err
	}

	hashPass, err := utils.HashPWD(data.Password)
	if err != nil {
		return 0, err
	}
	userToUpdate.ID = foundUser.ID
	userToUpdate.User_Type = foundUser.User_Type
	userToUpdate.Name = data.Name
	userToUpdate.Email = data.Email
	userToUpdate.Password = hashPass

	if data.Name == "" {
		userToUpdate.Name = foundUser.Name
	}
	if data.Email == "" {
		userToUpdate.Email = foundUser.Email
	}
	if data.Password == "" {
		userToUpdate.Password = foundUser.Password
	}
	update := bson.D{{Key: "$set", Value: userToUpdate}}
	result, err := s.userColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}
func (s *UserServicesImpl) DeleteUser(id string) (int64, error) {
	ctx, cansel := utils.Ctx()
	defer cansel()
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}
	filter := bson.D{{Key: "_id", Value: objId}}
	result, err := s.userColl.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
