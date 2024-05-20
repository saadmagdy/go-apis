package controllers

import (
	"basic_api/middleware"
	"basic_api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CartController struct {
	svc services.CartServices
}

func NewCartController(svc services.CartServices) CartController {
	return CartController{svc: svc}
}

func (crt *CartController) addProductsToCart(c *gin.Context) {
	uid := c.GetString("userId")
	pid := c.Param("id")
	err := crt.svc.AddProductToCart(uid, pid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "product add to cart successfully"})
}
func (crt *CartController) removeProductsFromCart(c *gin.Context) {
	uid := c.GetString("userId")
	pid := c.Param("id")
	err := crt.svc.RemoveProductFromCart(uid, pid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "product removed from cart successfully"})
}
func (crt *CartController) getUserCart(c *gin.Context) {
	uid := c.GetString("userId")

	cart, err := crt.svc.GetUserCart(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"cart": cart})

}

func (crt *CartController) InitRoutes(r *gin.Engine) {
	cart := r.Group("/cart")
	cart.Use(middleware.AuthMiddleware())

	cart.GET("/", middleware.AllowedTo("BUYER"), crt.getUserCart)
	cart.POST("/add/:id", middleware.AllowedTo("BUYER"), crt.addProductsToCart)
	cart.PATCH("/remove/:id", middleware.AllowedTo("BUYER"), crt.removeProductsFromCart)
}
