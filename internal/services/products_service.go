package services

import (
	"github.com/akhilnasimk/SS_backend/internal/models"
	"github.com/akhilnasimk/SS_backend/internal/repositories"
)

type ProductsService interface {
	GetAllProducts(limit int, offset int) ([]models.Product, int64 ,error)
}

type productsService struct {
	productRepo repositories.ProductsRepository
}

func NewProductsService(repo repositories.ProductsRepository) ProductsService {
	return &productsService{
		productRepo: repo,
	}
}

func (s *productsService) GetAllProducts(page, limit int) ([]models.Product, int64, error) {

    if limit <= 0 {
        limit = 10
    }
    if page <= 0 {
        page = 1
    }

    offset := (page - 1) * limit

    products, total, err := s.productRepo.GetAllProducts(limit, offset)
    if err != nil {
        return nil, 0, err
    }

    return products, total, nil
}
