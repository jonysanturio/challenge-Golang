package seeders

import (
  "log"
  "math/rand"
  "time"
  
  "gorm.io/gorm"
	"github.com/jonysanturio/challenge-golang/internal/models"
)
// RunSeeders runs all seeders
func RunSeeders(db *gorm.DB) error {
  rand.Seed(time.Now().UnixNano())
  // Seed categories
  if err := seedCategories(db); err != nil {
          return err
  }
  // Seed products
  if err := seedProducts(db); err != nil {
          return err
  }
  log.Println("Database seeding completed successfully")
  return nil
}
func seedCategories(db *gorm.DB) error {
  // Check if categories already exist
  var count int64
  db.Model(&models.Category{}).Count(&count)
  if count > 0 {
          log.Println("Categories already seeded, skipping")
          return nil
  }
  categories := []models.Category{
          {Name: "Electronics", Description: "Electronic devices and gadgets"},
          {Name: "Clothing", Description: "Apparel and fashion items"},
          {Name: "Books", Description: "Books and educational materials"},
          {Name: "Home & Garden", Description: "Home appliances and garden supplies"},
          {Name: "Sports", Description: "Sports equipment and accessories"},
  }
  for _, category := range categories {
          if err := db.Create(&category).Error; err != nil {
                  return err
          }
  }
  log.Printf("Seeded %d categories", len(categories))
  return nil
}
func seedProducts(db *gorm.DB) error {
  // Check if products already exist
  var count int64
  db.Model(&models.Product{}).Count(&count)
  if count > 0 {
          log.Println("Products already seeded, skipping")
          return nil
  }
  // Get all categories for assignment
  var categories []models.Category
  if err := db.Find(&categories).Error; err != nil {
          return err
  }
  if len(categories) == 0 {
          log.Println("No categories found for product seeding")
          return nil
  }
  products := []models.Product{
          {Name: "Smartphone X", Description: "Latest flagship smartphone", Price: 999.99, Stock: 50},
          {Name: "Laptop Pro", Description: "High-performance laptop for professionals", Price: 1499.99, Stock: 30},
          {Name: "Wireless Headphones", Description: "Noise-cancelling wireless headphones", Price: 299.99, Stock: 100},
          {Name: "Cotton T-Shirt", Description: "Comfortable 100% cotton t-shirt", Price: 19.99, Stock: 200},
          {Name: "Denim Jeans", Description: "Classic fit denim jeans", Price: 79.99, Stock: 150},
          {Name: "Running Shoes", Description: "Lightweight running shoes for athletes", Price: 129.99, Stock: 75},
          {Name: "Programming Book", Description: "Learn to code with this comprehensive guide", Price: 49.99, Stock: 80},
          {Name: "Coffee Maker", Description: "Automatic drip coffee maker", Price: 89.99, Stock: 60},
          {Name: "Garden Hose", Description: "50-foot expandable garden hose", Price: 34.99, Stock: 90},
          {Name: "Yoga Mat", Description: "Non-slip yoga mat for practice", Price: 24.99, Stock: 120},
  }
  for _, product := range products {
          if err := db.Create(&product).Error; err != nil {
                  return err
          }
          // Assign 1-3 random categories to each product
          numCategories := 1 + rand.Intn(3) // 1 to 3 categories
          selectedCategories := make([]models.Category, 0, numCategories)
          // Simple random selection (could be improved)
          perm := rand.Perm(len(categories))
          for j := 0; j < numCategories && j < len(perm); j++ {
                  selectedCategories = append(selectedCategories, categories[perm[j]])
          }
          if err := db.Model(&product).Association("Categories").Append(&selectedCategories); err != nil {
                  return err
          }
          // Create initial history record
          history := models.ProductHistory{
                  ProductID: product.ID,
                  Price:     product.Price,
                  Stock:     product.Stock,
                  ChangedAt: time.Now(),
          }
          if err := db.Create(&history).Error; err != nil {
                  return err
          }
  }
  log.Printf("Seeded %d products with categories and history", len(products))
  return nil
}