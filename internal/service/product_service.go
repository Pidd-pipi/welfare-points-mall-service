package service

import (
	"fmt"
	"log/slog"

	"gorm.io/gorm"

	"github.com/ld/welfaremall/internal/constants"
	"github.com/ld/welfaremall/internal/model"
	"github.com/ld/welfaremall/internal/repository"
	"github.com/ld/welfaremall/internal/util"
)

// ProductService 商品业务逻辑。
type ProductService interface {
	Create(name string, category constants.ProductCategory, pointsCost, stock, exchangeLimit int, coverImage, description string) (*model.Product, error)
	Update(id uint, name string, category constants.ProductCategory, pointsCost, stock, exchangeLimit int, coverImage, description string) (*model.Product, error)
	Delete(id uint) error
	ToggleStatus(id uint) (*model.Product, error)
	List(page, pageSize int, category constants.ProductCategory, keyword string) ([]model.Product, int64, error)
	GetByID(id uint) (*model.Product, error)
	DecrementStock(id uint, qty int) error
	DecrementStockTx(tx *gorm.DB, id uint, qty int) error
	IncrementStock(id uint, qty int) error
	IncrementStockTx(tx *gorm.DB, id uint, qty int) error
}

type productService struct {
	productRepo repository.ProductRepository
	db          *gorm.DB
	logger      *slog.Logger
}

// NewProductService 构造商品服务。
func NewProductService(productRepo repository.ProductRepository, db *gorm.DB, logger *slog.Logger) ProductService {
	return &productService{productRepo: productRepo, db: db, logger: logger}
}

func (s *productService) Create(name string, category constants.ProductCategory, pointsCost, stock, exchangeLimit int, coverImage, description string) (*model.Product, error) {
	if name == "" || pointsCost <= 0 {
		return nil, fmt.Errorf("create product: %w", util.ErrValidation)
	}
	if !category.Valid() {
		category = constants.CategoryPhysical
	}
	if exchangeLimit <= 0 {
		exchangeLimit = 1
	}
	product := &model.Product{Name: name, Category: category, PointsCost: pointsCost, Stock: stock, ExchangeLimit: exchangeLimit, Status: "on", CoverImage: coverImage, Description: description}
	if err := s.productRepo.Create(product); err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}
	s.logger.Info(constants.LogProductCreated, "product_id", product.ID, "name", name)
	return product, nil
}

func (s *productService) Update(id uint, name string, category constants.ProductCategory, pointsCost, stock, exchangeLimit int, coverImage, description string) (*model.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("update product[id=%d]: %w", id, err)
	}
	if name != "" {
		product.Name = name
	}
	if category.Valid() {
		product.Category = category
	}
	if pointsCost > 0 {
		product.PointsCost = pointsCost
	}
	if stock >= 0 {
		product.Stock = stock
	}
	if exchangeLimit > 0 {
		product.ExchangeLimit = exchangeLimit
	}
	if coverImage != "" {
		product.CoverImage = coverImage
	}
	if description != "" {
		product.Description = description
	}
	if err := s.productRepo.Update(product); err != nil {
		return nil, fmt.Errorf("update product[id=%d]: %w", id, err)
	}
	s.logger.Info(constants.LogProductUpdated, "product_id", id)
	return product, nil
}

func (s *productService) Delete(id uint) error {
	if err := s.productRepo.Delete(id); err != nil {
		return fmt.Errorf("delete product[id=%d]: %w", id, err)
	}
	s.logger.Info(constants.LogProductDeleted, "product_id", id)
	return nil
}

func (s *productService) ToggleStatus(id uint) (*model.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("toggle product[id=%d]: %w", id, err)
	}
	if product.Status == "on" {
		product.Status = "off"
	} else {
		product.Status = "on"
	}
	if err := s.productRepo.Update(product); err != nil {
		return nil, fmt.Errorf("toggle product[id=%d]: %w", id, err)
	}
	s.logger.Info(constants.LogProductStatusToggled, "product_id", id, "status", product.Status)
	return product, nil
}

func (s *productService) List(page, pageSize int, category constants.ProductCategory, keyword string) ([]model.Product, int64, error) {
	products, total, err := s.productRepo.List(page, pageSize, category, keyword)
	if err != nil {
		return nil, 0, fmt.Errorf("list products: %w", err)
	}
	s.logger.Info(constants.LogProductListQueried, "total", total)
	return products, total, nil
}

func (s *productService) GetByID(id uint) (*model.Product, error) {
	return s.productRepo.FindByID(id)
}

func (s *productService) DecrementStock(id uint, qty int) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return s.DecrementStockTx(tx, id, qty)
	})
}

func (s *productService) DecrementStockTx(tx *gorm.DB, id uint, qty int) error {
	if qty <= 0 {
		return fmt.Errorf("decrement stock product[id=%d] qty[%d]: %w", id, qty, util.ErrValidation)
	}
	if err := s.productRepo.IncrementStockTx(tx, id, qty); err != nil {
		return fmt.Errorf("decrement stock product[id=%d]: %w", id, err)
	}
	return nil
}

func (s *productService) IncrementStock(id uint, qty int) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return s.IncrementStockTx(tx, id, qty)
	})
}

func (s *productService) IncrementStockTx(tx *gorm.DB, id uint, qty int) error {
	if qty <= 0 {
		return fmt.Errorf("increment stock product[id=%d] qty[%d]: %w", id, qty, util.ErrValidation)
	}
	if err := s.productRepo.DecrementStockTx(tx, id, qty); err != nil {
		return fmt.Errorf("increment stock product[id=%d]: %w", id, err)
	}
	return nil
}
