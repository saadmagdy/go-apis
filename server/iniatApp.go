package server

import (
	"basic_api/controllers"
	"basic_api/database"
	"basic_api/services"
	"basic_api/utils"
	"log"
	"os"

	"github.com/lpernett/godotenv"
)

var (
	userCtrl    controllers.UserController
	productCtrl controllers.ProductController
	cartCtrl    controllers.CartController
)

func INIATAPP() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// ************************************************************************************************
	port := os.Getenv("PORT")
	if port == "" {
		port = ":5000" // Default Port if not set in environment variable PORT
	}
	dbName := os.Getenv("DB_NAME")
	dbUri := os.Getenv("DB_URI")
	// ************************************************************************************************
	ctx, cansel := utils.Ctx()
	defer cansel()
	//************************************************************************************************
	r := NewServer(":" + port)
	//************************************************************************************************
	mongodb := database.NewMongoDb(ctx, dbName, dbUri)
	client, err := mongodb.CreateClient()
	if err != nil {
		log.Fatal(err)
	}
	// ************************************************************************************************
	defer func() {
		err = mongodb.CloseConnection(client)
		if err != nil {
			log.Fatal(err)
		}
	}()
	//*********************************************************************************************
	db := mongodb.CreateDataBase(client)
	// ************************************************************************************************
	uCollName := os.Getenv("USERS_COLL")
	userColl := mongodb.CreateCollection(db, uCollName)
	usersvc := services.NewUserService(userColl)
	userCtrl = controllers.NewUserController(usersvc)
	//************************************************************************************************
	pCollName := os.Getenv("PRODUCTS_COLL")
	pColl := mongodb.CreateCollection(db, pCollName)
	psvc := services.NewProductServices(pColl, userColl)
	productCtrl = controllers.NewProductCtrl(psvc)
	// ************************************************************************************************
	cartSCollName := os.Getenv("CARTS_COLL")
	cartsColl := mongodb.CreateCollection(db, cartSCollName)
	cartServices := services.NewCartServices(cartsColl, pColl)
	cartCtrl = controllers.NewCartController(cartServices)
	//**********************************************************************************************
	if err := r.Run(); err != nil {
		log.Fatal("Error creatring Server")
	}
	// ************************************************************************************************
}
