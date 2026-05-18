package handlers

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"

	"github.com/jonysanturio/challenge-golang/config"
	"github.com/jonysanturio/challenge-golang/internal/models"
	"github.com/jonysanturio/challenge-golang/internal/services"
)

type ProductHandler struct {
	service services.ProductService
}

func NewProductHandler(service services.ProductService) *ProductHandler {
	return &ProductHandler{service}
}

func (h *ProductHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var input models.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateProduct(uint(id), &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar producto"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Producto actualizado"})
}

// Implementacion de los handlers
func GetProducts(c *gin.Context) {
	db := config.GetDB()
	var product []models.Product
	if err := db.Find(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener productos"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func GetProductByID(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}

func GetProduct(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}