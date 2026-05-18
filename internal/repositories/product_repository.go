package repositories

import (
	"github.com/jonysanturio/challenge-golang/internal/models"

	"gorm.io/gorm"
)

type ProductRepository interface{
	UpdateWithHistory(product *models.Product, history *models.ProductHistory) error
}

type productRepository struct{
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository{
	return &productRepository{db}
}

func (r *productRepository) UpdateWithHistory(product *models.Product, history *models.ProductHistory) error{
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(product).Error; err != nil {
			return err
		}

		if history.Price == 0 && history.Stock == 0 {
			return nil
		}
		if err := tx.Create(history).Error; err != nil {
			return err
		}
		return nil
	})
}