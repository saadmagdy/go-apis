package controllers

import (
	"basic_api/middleware"
	"basic_api/models"
	"basic_api/services"
	"basic_api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserController struct {
	svc services.UserServices
}

var validate = validator.New()

func NewUserController(svc services.UserServices) UserController {
	return UserController{svc: svc}
}
func (u *UserController) getAllUsers(c *gin.Context) {
	users, err := u.svc.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(users) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No users found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}
func (u *UserController) getUserById(c *gin.Context) {
	id := c.Param("id")
	user, err := u.svc.GetUserById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, user)
}
func (u *UserController) signUp(c *gin.Context) {
	var user *models.UserSignUp
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	validationErr := validate.Struct(user)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}
	userExist, _ := u.svc.GetUserByEmail(user.Email)
	if userExist != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email has been used."})
		return
	}

	if err := u.svc.CreateUser(user); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error(), "msg": "error creating user"})
		return

	}
	c.JSON(200, gin.H{"message": "success creating the user"})
}
func (u *UserController) login(c *gin.Context) {
	var loginuser *models.UserLogin
	if err := c.BindJSON(&loginuser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if validationErr := validate.Struct(loginuser); validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}
	user, err := u.svc.GetUserByEmail(loginuser.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect email address or password"})
		return
	}
	k, msg := utils.VerfiyPWD(loginuser.Password, user.Password)
	if !k {
		c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
		return
	}
	token, err := user.GenerateToken()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "token_error": "token error"})
		return
	}
	user.Password = ""
	c.JSON(200, gin.H{"message": "LogIn successfully", "token": token, "user": user})
}
func (u *UserController) updateUser(c *gin.Context) {
	id := c.GetString("userId")
	var data *models.UserUpdate
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	i, err := u.svc.UpdateUser(id, data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if i < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error updating user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success Updating the user"})
}
func (u *UserController) deleteUser(c *gin.Context) {
	id := c.GetString("userId")
	i, err := u.svc.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error(), "id": id})
		return
	}
	if i < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error Deleting Document"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted user successfully"})
}
func (u *UserController) InitRoutes(r *gin.Engine) {
	users := r.Group("/users")
	users.POST("/signup", u.signUp)
	users.POST("/login", u.login)
	users.Use(middleware.AuthMiddleware())

	users.GET("/", u.getAllUsers)
	users.PATCH("/", u.updateUser)
	users.DELETE("/", u.deleteUser)
	users.GET("/:id", u.getUserById)
}
