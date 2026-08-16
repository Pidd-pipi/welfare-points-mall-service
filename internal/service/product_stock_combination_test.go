package service

import (
	"log/slog"
	"os"
	"testing"

	"gorm.io/gorm"

	"github.com/ld/welfaremall/internal/constants"
	"github.com/ld/welfaremall/internal/model"
	"github.com/ld/welfaremall/internal/util"
)

type mockProductRepo struct {
	stock map[uint]int
}

func (m *mockProductRepo) Create(product *model.Product) error { return nil }
func (m *mockProductRepo) FindByID(id uint) (*model.Product, error) {
	return &model.Product{ID: id, Stock: m.stock[id], PointsCost: 100}, nil
}
func (m *mockProductRepo) FindByIDTx(tx *gorm.DB, id uint) (*model.Product, error) {
	return m.FindByID(id)
}
func (m *mockProductRepo) List(page, pageSize int, category constants.ProductCategory, keyword string) ([]model.Product, int64, error) {
	return nil, 0, nil
}
func (m *mockProductRepo) Update(product *model.Product) error { return nil }
func (m *mockProductRepo) UpdateTx(tx *gorm.DB, product *model.Product) error { return nil }
func (m *mockProductRepo) DecrementStockTx(tx *gorm.DB, id uint, qty int) error {
	m.stock[id] -= qty
	if m.stock[id] < 0 {
		return util.ErrProductOutOfStock
	}
	return nil
}
func (m *mockProductRepo) IncrementStockTx(tx *gorm.DB, id uint, qty int) error {
	m.stock[id] += qty
	return nil
}
func (m *mockProductRepo) Delete(id uint) error { return nil }

func newProductSvc() ProductService {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	return NewProductService(&mockProductRepo{stock: map[uint]int{1: 10}}, nil, logger)
}

func TestProductDecrementStock(t *testing.T) {
	svc := newProductSvc()
	if err := svc.DecrementStockTx(nil, 1, 3); err != nil {
		t.Fatalf("decrement failed: %v", err)
	}
	product, _ := svc.GetByID(1)
	if product.Stock != 7 {
		t.Fatalf("stock = %d, want 7", product.Stock)
	}
}

func TestExchangeCostBoundary(t *testing.T) {
	if got := util.ExchangeCost(100, 3); got != 300 {
		t.Fatalf("ExchangeCost(100,3) = %d, want 300", got)
	}
}
