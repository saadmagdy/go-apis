package controllers

import (
	"basic_api/middleware"
	"basic_api/models"
	"basic_api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	svc services.ProductServices
}

func NewProductCtrl(svc services.ProductServices) ProductController {
	return ProductController{
		svc: svc,
	}
}

func (p *ProductController) createProduct(c *gin.Context) {
	uid := c.GetString("userId")
	var data *models.ProductCreate

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if verr := validate.Struct(data); verr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": verr.Error()})
		return
	}
	if err := p.svc.CreateProduct(data, uid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"msg": "Product created successfully"})
}
func (p *ProductController) getProductById(c *gin.Context) {
	pid := c.Param("id")
	if pid == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "product id requierd"})
		return
	}
	prod, err := p.svc.GetProductById(pid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"msg": "success getting the product", "prouct": prod})
}
func (p *ProductController) updateProductById(c *gin.Context) {
	pid := c.Param("id")
	if pid == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "product id is required"})
		return
	}
	uid := c.GetString("userId")
	if uid == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var data *models.ProductUpdate
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uc, err := p.svc.UpdateProduct(data, pid, uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else if uc < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error updating product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "product updated successfully"})
}
func (p *ProductController) deleteProductById(c *gin.Context) {
	pid := c.Param("id")
	if pid == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "product id is required"})
		return
	}
	uid := c.GetString("userId")
	if uid == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	dc, err := p.svc.DeleteProduct(pid, uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else if dc < 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Product not found"})
	}
	c.JSON(http.StatusOK, gin.H{"msg": "product deleted successfully"})
}
func (p *ProductController) getAllProducts(c *gin.Context) {
	prods, err := p.svc.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "success getting all products", "products": prods})
}

func (p *ProductController) InitRoutes(r *gin.Engine) {
	productGroup := r.Group("/products")
	productGroup.Use(middleware.AuthMiddleware())
	productGroup.POST("/", middleware.AllowedTo("SELLER"), p.createProduct)
	productGroup.PATCH("/:id", middleware.AllowedTo("SELLER"), p.updateProductById)
	productGroup.DELETE("/:id", middleware.AllowedTo("SELLER"), p.deleteProductById)

	productGroup.GET("/:id", p.getProductById)
	productGroup.GET("/", p.getAllProducts)
}
