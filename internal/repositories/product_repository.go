package repositories

import (
	"challenge-golang/internal/models"

	"github.com/jonysanturio/qisur-challenge/internal/models"
	"gorm.io/gorm"
)

type ProductRepository interface{
	UpdateWithHistory(product *models.Product, history *models.Product) error
}

type productRepository struct{
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository{
	return &productRepository{db}
}

func (r *productRepository) UpdateWithHistory(product *models.Product, history *models.Product) error{
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Actualizacion del producto
		if err := tx.Save(product).Error; err != nil {
			return err
		}

		// Guardado de historial
		if history.Price == 0 && history.Stock == 0 {
			return nil
		}
		if err := tx.Create(history).Error; err != nil {
			return err
		}
		return nil
	})
}