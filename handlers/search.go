package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"

    "qisur-challenge/models"
    "qisur-challenge/config"
    "gorm.io/gorm"
)
   
// Search handles GET /api/search
// Supports searching products or categories with query parameters
func Search(c *gin.Context) {
    db := config.GetDB()

// Search type: product or category
searchType := c.Query("type")
if searchType == "" {
        searchType = "product" // default
}
// Pagination
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
offset := (page - 1) * lim
// Search query
query := c.Query("q")
if query == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Search que required"})
        return
}
   
switch searchType {
case "product":
        searchProducts(c, db, query, page, limit, offset)
case "category":
        searchCategories(c, db, query, page, limit, offset)
default:
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid setype. Use 'product' or 'category'"})
}
}

func searchProducts(c *gin.Context, db *gorm.DB, query string,
 page, limit,
 offset int) {
  var products []models.Product
  var total int64

  db.Model(&models.Product{}).
          Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%"+query+"%").
          Count(&total).

  Preload("Categories").
          Limit(limit).
          Offset(offset).
          Order("created_at DESC").
          Find(&products)

  c.JSON(http.StatusOK, gin.H{
          "data": products,
          "pagination": gin.H{
                  "page":  page,
                  "limit": limit,
                  "total": total,
                  "pages": (total + int64(limit) - 1) / int64(limit)

          },
          "search": gin.H{
                  "type": "product",
                  "query": query,
          },
  })
}

func searchCategories(c *gin.Context, db *gorm.DB, query string, page, limit, offset int) {
  var categories []models.Category
  var total int64

  db.Model(&models.Category{}).
        Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%"+query+"%").Count(&total).

Preload("Products").
          Limit(limit).
          Offset(offset).
          Order("created_at DESC").
          Find(&categories)

c.JSON(http.StatusOK, gin.H{
        "data": categories,
        "pagination": gin.H{
                "page":  page,
                "limit": limit,
                "total": total,
                "pages": (total + int64(limit) - 1) / int64(limit),
        },
        "search": gin.H{
                "type": "category",
                "query": query,
        },
})
}