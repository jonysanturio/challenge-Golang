package models

import (
	"time"

	"gorm.io/gorm"
)

// Modelo de producto
type Product struct{
	ID		uint	`gorm:"primaryKey" json:"id"`
	Name   	string 	`gorm:"size:255;not null" json:"name" binding:"required"`
	Description string `gorm:"type:text" json:"description"`
	Price 	float64  `gorm:"type:decimal(10,2);not null" json:"price" binding:"required"`
	Stock	int		`gorm:"not null" json:"stock" binding:"required",gte=0`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Categories []Category `gorm:"many2many:product_categories;" json:"categories,omitempty"`
}

// Producto Historico
type ProductHistory struct {
	ID		uint	`gorm:"primaryKey" json:"id"`
	ProductID uint `gorm:"not null;index" json:"product_id"`
	Price	float64	`gorm:"type:decimal(10,2);not null" json:"price"`
	Stock	int		`gorm:"not null" json:"stock"`
	ChangedAt	time.Time	`gorm:"not null;default:now()" json:"changed"`
	
	Product Product `gorm:"foreignKey:ProductID;constraint:OneDelete:CASCADE;" json:"-"`
}

// Categorias
type Category struct{
	ID		uint	`gorm:"primaryKey" json:"id"`
	Name	string	`gorm:"size:255;not null;unique" json:"name" binding:"required"`
	Description string `gorm:"type:text" json:"description"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
	DeletedAt	gorm.DeletedAt `gorm:"index" json:"-"`

	Products []Product `gorm:"many2many:product_categories;" json:"products,omitempty"`
}

// Tabla intermedia para la relacion muchos a muchos
type ProductCategory struct{
	ProductID uint `gorm:"primaryKey"`
	CategoryID uint `gorm:"primaryKey"`
}