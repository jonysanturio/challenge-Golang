package handlers

import (
	"net/http"
	"strconv"
	"github.com/jonysanturio/challenge-golang/internal/models"
	"github.com/jonysanturio/challenge-golang/internal/services"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service services.ProductService
}

func NewProductHandler(service services.ProductService) *ProductHandler {
	return &ProductHandler{service}
}

func (h *ProductHandler) Update(c *gin.Context) {
	// 1. Extraer ID
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// 2. Parsear y validar JSON
	var input models.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Llamar al servicio (Limpio y delegando responsabilidad)
	if err := h.service.UpdateProduct(uint(id), &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar producto"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Producto actualizado y evento emitido"})
}

// Dummy implementations for missing handlers
func GetProducts(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
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