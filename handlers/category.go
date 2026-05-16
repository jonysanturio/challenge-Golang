package handlers

import (
  "net/http"
  "strconv"

  "github.com/gin-gonic/gin"

  "qisur-challenge/models"
  "qisur-challenge/config"
)

// GetCategories handles GET /api/categories
func GetCategories(c *gin.Context) {
  db := config.GetDB()

  // Pagination parameters
  page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
  if err != nil || page < 1 {
          page = 1
  }
  limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
  if err != nil || limit < 1 {
          limit = 10
  }
  if limit > 100 {
          limit = 100
  }
  offset := (page - 1) * limit

  // Filtering
  name := c.Query("name")

  var categories []models.Category
  var total int64

  query := db.Model(&models.Category{})

  if name != "" {
          query = query.Where("name ILIKE ?", "%"+name+"%")
  }

  query.Count(&total)
  query = db.Offset(limit).Limit(limit).Offset(offset).Find(&cate
gories)

  c.JSON(http.StatusOK, gin.H{
          "data": categories,
          "pagination": gin.H{
                  "page":  page,
                  "limit": limit,
                  "total": total,
                  "pages": (total + int64(limit) - 1) / int64(limi
          },
  })
}

// GetCategoryByID handles GET /api/categories/:id
func GetCategoryByID(c *gin.Context) {
  db := config.GetDB()
  id := c.Param("id")

  var category models.Category
  if err := db.Preload("Products").First(&category, id).Error; err != nil {
          c.JSON(http.StatusNotFound, gin.H{"error": "Category not"})
          return
  }

  c.JSON(http.StatusOK, gin.H{"data": category})
}

// CreateCategory handles POST /api/categories
func CreateCategory(c *gin.Context) {
  db := config.GetDB()
  var input struct {
          Name        string `json:"name" binding:"required"`
          Description string `json:"description"`
  }

  if err := c.ShouldBindJSON(&input); err != nil {
          c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
          return
  }

  category := models.Category{
          Name:        input.Name,
          Description: input.Description,
  }

  if err := db.Create(&category).Error; err != nil {
          c.JSON(http.StatusInternalServerError, gin.H{"error": "For create category"})
          return
  }

  // Notify via WebSocket
  // websocket.NotifyCategoryCreated(category)

  c.JSON(http.StatusCreated, gin.H{"data": category})
}

// UpdateCategory handles PUT /api/categories/:id
func UpdateCategory(c *gin.Context) {
  db := config.GetDB()
  id := c.Param("id")

  var category models.Category
  if err := db.First(&category, id).Error; err != nil {
          c.JSON(http.StatusNotFound, gin.H{"error": "Category not"})
          return
  }

  var input struct {
          Name        string `json:"name"`
          Description string `json:"description"`
  }

  if err := c.ShouldBindJSON(&input); err != nil {
          c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
          return
  }

  if input.Name != "" {
          category.Name = input.Name
  }
  if input.Description != "" {
          category.Description = input.Description
  }

  if err := db.Save(&category).Error; err != nil {
          c.JSON(http.StatusInternalServerError, gin.H{"error": "For update category"})
          return
  }

  // Notify via WebSocket
  // websocket.NotifyCategoryUpdated(category)

  c.JSON(http.StatusOK, gin.H{"data": category})
}

// DeleteCategory handles DELETE /api/categories/:id
func DeleteCategory(c *gin.Context) {
  db := config.GetDB()
  id := c.Param("id")

  var category models.Category
  if err := db.First(&category, id).Error; err != nil {
          c.JSON(http.StatusNotFound, gin.H{"error": "Category not"})
          return
  }

  if err := db.Delete(&category).Error; err != nil {
          c.JSON(http.StatusInternalServerError, gin.H{"error": "For delete category"})
          return
  }

  // Notify via WebSocket
  // websocket.NotifyCategoryDeleted(category.ID)

  c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
