package services

import (
	"time"

	"github.com/jonysanturio/challenge-golang/internal/models"
	"github.com/jonysanturio/challenge-golang/internal/repositories"
	"github.com/jonysanturio/challenge-golang/websocket"

)

type ProductService interface {
	UpdateProduct(id uint, input *models.Product) error	
}

type productService struct{
	repo repositories.ProductRepository
	hub *websocket.Hub
}

func NewProductService(repo repositories.ProductRepository, hub *websocket.Hub) ProductService{
	return &productService{repo: repo, hub: hub}
}

func (s *productService) UpdateProduct(id uint, input *models.Product) error {
	input.ID = id

	history := &models.ProductHistory{
		ProductID:    id,
		Stock:        input.Stock,
		Price:        input.Price,
		ChangedAt:    time.Now(),
	}
	err := s.repo.UpdateWithHistory(input, history)
	if err != nil{
		return err
	}
	eventMsg := []byte(`{"event": "PRODUCT_UPDATED", "product_id": ` + string(rune(id)) + `}`)
	s.hub.Broadcast <- eventMsg
	
	return nil
}
