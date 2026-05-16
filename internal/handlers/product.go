package handlers

import (
   "net/http"
   "strconv"
   "time"

   "github.com/gin-gonic/gin"

   "qisur-challenge/models"
   "qisur-challenge/config"
 )

 // GetProducts handles GET /api/products
 // Returns paginated list of products with optional filtering and sorting
 func GetProducts(c *gin.Context) {
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
   minPrice, _ := strconv.ParseFloat(c.Query("min_price"), 64)
   maxPrice, _ := strconv.ParseFloat(c.Query("max_price"), 64)
   minStock, _ := strconv.Atoi(c.Query("min_stock"))
   maxStock, _ := strconv.Atoi(c.Query("max_stock"))

   // Sorting
   sortBy := c.DefaultQuery("sort_by", "created_at")
   sortOrder := c.DefaultQuery("sort_order", "desc")
   if sortOrder != "asc" && sortOrder != "desc" {
           sortOrder = "desc"
   }

   var products []models.Product
   var total int64

   query := db.Model(&models.Product{})

   // Apply filters
   if name != "" {
           query = query.Where("name ILIKE ?", "%"+name+"%")
   }
   if minPrice > 0 {
           query = query.Where("price >= ?", minPrice)
   }
   if maxPrice > 0 {
           query = query.Where("price <= ?", maxPrice)
   }
   if minStock >= 0 {
           query = query.Where("stock >= ?", minStock)
   }
   if maxStock >= 0 {
           query = query.Where("stock <= ?", maxStock)
   }

   // Get total count
   query.Count(&total)

   // Apply sorting and pagination
   query = db.Order(sortBy + " " + sortOrder).
           Limit(limit).
           Offset(offset).
           Find(&products)

   c.JSON(http.StatusOK, gin.H{
           "data":      products,
           "pagination": gin.H{
                   "page":  page,
                   "limit": limit,
                   "total": total,
                   "pages": (total + int64(limit) - 1) / int64(limit),
           },
   })
 }

 // GetProductByID handles GET /api/products/:id
 func GetProductByID(c *gin.Context) {
   db := config.GetDB()
   id := c.Param("id")

   var product models.Product
   if err := db.Preload("Categories").First(&product, id).Error; e
 rr != nil {
           c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"}
 )
           return
   }

   c.JSON(http.StatusOK, gin.H{"data": product})
 }

 // CreateProduct handles POST /api/products
 func CreateProduct(c *gin.Context) {
   db := config.GetDB()

   var input struct {
           Name        string  `json:"name" binding:"required"`
           Description string  `json:"description"`
           Price       float64 `json:"price" binding:"required,gt=0"`
           Stock       int     `json:"stock" binding:"required,gte=0"`
           CategoryIDs []uint  `json:"category_ids"`
   }

   if err := c.ShouldBindJSON(&input); err != nil {
           c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
  }

  // Create product
  product := models.Product{
          Name:        input.Name,
          Description: input.Description,
          Price:       input.Price,
          Stock:       input.Stock,
  }

  if err := db.Create(&product).Error; err != nil {
          c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
          return
  }

  // Assign categories if provided
  if len(input.CategoryIDs) > 0 {
          var categories []models.Category
          if err := db.Find(&categories,
  input.CategoryIDs).Error; err != nil {
                   c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category IDs"})
                   return
           }

           db.Model(&product).Association("Categories").Replace(&categories)
   }

   // Create history record
   history := models.ProductHistory{
           ProductID: product.ID,
           Price:     product.Price,
           Stock:     product.Stock,
           ChangedAt: time.Now(),
   }
   db.Create(&history)

   // Reload with associations
   db.Preload("Categories").First(&product, product.ID)

   // Notify via WebSocket (to be implemented in websocket package)
   // websocket.NotifyProductCreated(product)

   c.JSON(http.StatusCreated, gin.H{"data": product})
 }

 // UpdateProduct handles PUT /api/products/:id
 func UpdateProduct(c *gin.Context) {
   db := config.GetDB()
   id := c.Param("id")

   var product models.Product
   if err := db.First(&product, id).Error; err != nil {
           c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"}

 )
           return
   }

   var input struct {
           Name        string  `json:"name"`
           Description string  `json:"description"`
           Price       float64 `json:"price"`
           Stock       int     `json:"stock"`
           CategoryIDs []uint  `json:"category_ids"`
   }

   if err := c.ShouldBindJSON(&input); err != nil {
           c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
           return
   }

   // Update fields if provided
   if input.Name != "" {
           product.Name = input.Name
   }
   if input.Description != "" {
           product.Description = input.Description
   }
   if input.Price > 0 {
           product.Price = input.Price
   }
   if input.Stock >= 0 {
           product.Stock = input.Stock
   }

   if err := db.Save(&product).Error; err != nil {
           c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
           return
   }

   // Update categories if provided
   if input.CategoryIDs != nil {
           var categories []models.Category
           if err := db.Find(&categories, input.CategoryIDs).Error; err != nil {
                   c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category IDs"})
                   return
           }

           db.Model(&product).Association("Categories").Replace(&categories)
   }

   // Create history record
   history := models.ProductHistory{
           ProductID: product.ID,
           Price:     product.Price,
           Stock:     product.Stock,
           ChangedAt: time.Now(),
   }
   db.Create(&history)

   // Reload with associations
   db.Preload("Categories").First(&product, product.ID)

   // Notify via WebSocket
   // websocket.NotifyProductUpdated(product)

   c.JSON(http.StatusOK, gin.H{"data": product})
 }

 // DeleteProduct handles DELETE /api/products/:id
 func DeleteProduct(c *gin.Context) {
   db := config.GetDB()
   id := c.Param("id")

   var product models.Product
   if err := db.First(&product, id).Error; err != nil {
           c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
           return
   }

   if err := db.Delete(&product).Error; err != nil {
           c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
           return
   }

   // Notify via WebSocket
   // websocket.NotifyProductDeleted(product.ID)

   c.JSON(http.StatusOK, gin.H{"message": "Product deleted success fully"})
 }

 // GetProductHistory handles GET /api/products/:id/history
 func GetProductHistory(c *gin.Context) {
   db := config.GetDB()
   id := c.Param("id")

   // Validate product exists
   var product models.Product
   if err := db.First(&product, id).Error; err != nil {
           c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"}
 )
           return
   }

   // Date range filtering
   startDate := c.Query("start")
   endDate := c.Query("end")

   var history []models.ProductHistory
   query := db.Model(&models.ProductHistory{}).
           Where("product_id = ?", id)

   if startDate != "" {
           query = query.Where("changed_at >= ?", startDate)
   }
   if endDate != "" {
           query = query.Where("changed_at <= ?", endDate)
   }

   query.Order("changed_at DESC").Find(&history)

   c.JSON(http.StatusOK, gin.H{"data": history})
 }