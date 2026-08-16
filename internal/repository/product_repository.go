package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/ld/welfaremall/internal/constants"
	"github.com/ld/welfaremall/internal/model"
	"github.com/ld/welfaremall/internal/util"
)

// ProductRepository 商品仓储。
type ProductRepository interface {
	Create(product *model.Product) error
	FindByID(id uint) (*model.Product, error)
	FindByIDTx(tx *gorm.DB, id uint) (*model.Product, error)
	List(page, pageSize int, category constants.ProductCategory, keyword string) ([]model.Product, int64, error)
	Update(product *model.Product) error
	UpdateTx(tx *gorm.DB, product *model.Product) error
	DecrementStockTx(tx *gorm.DB, id uint, qty int) error
	IncrementStockTx(tx *gorm.DB, id uint, qty int) error
	Delete(id uint) error
}

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository 构造商品仓储。
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *model.Product) error {
	if err := r.db.Create(product).Error; err != nil {
		return fmt.Errorf("create product: %w", err)
	}
	return nil
}

func (r *productRepository) FindByID(id uint) (*model.Product, error) {
	return r.FindByIDTx(nil, id)
}

func (r *productRepository) FindByIDTx(tx *gorm.DB, id uint) (*model.Product, error) {
	var p model.Product
	err := dbOrTx(r.db, tx).First(&p, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("find product by id: %w", util.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("find product by id: %w", err)
	}
	return &p, nil
}

func (r *productRepository) List(page, pageSize int, category constants.ProductCategory, keyword string) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64
	q := r.db.Model(&model.Product{})
	if category != "" {
		q = q.Where("category = ?", category)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("name LIKE ? OR description LIKE ?", like, like)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count products: %w", err)
	}
	if err := q.Offset(page * pageSize).Limit(pageSize).Order("id desc").Find(&products).Error; err != nil {
		return nil, 0, fmt.Errorf("list products: %w", err)
	}
	return products, total, nil
}

func (r *productRepository) Update(product *model.Product) error {
	return r.UpdateTx(nil, product)
}

func (r *productRepository) UpdateTx(tx *gorm.DB, product *model.Product) error {
	if err := dbOrTx(r.db, tx).Save(product).Error; err != nil {
		return fmt.Errorf("update product: %w", err)
	}
	return nil
}

func (r *productRepository) DecrementStockTx(tx *gorm.DB, id uint, qty int) error {
	res := dbOrTx(r.db, tx).Model(&model.Product{}).
		Where("id = ? AND stock >= ?", id, qty).
		UpdateColumn("stock", gorm.Expr("stock - ?", qty))
	if res.Error != nil {
		return fmt.Errorf("decrement product stock: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("decrement product stock: %w", util.ErrProductOutOfStock)
	}
	return nil
}

func (r *productRepository) IncrementStockTx(tx *gorm.DB, id uint, qty int) error {
	res := dbOrTx(r.db, tx).Model(&model.Product{}).
		Where("id = ?", id).
		UpdateColumn("stock", gorm.Expr("stock + ?", qty))
	if res.Error != nil {
		return fmt.Errorf("increment product stock: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("increment product stock: %w", util.ErrNotFound)
	}
	return nil
}

func (r *productRepository) Delete(id uint) error {
	res := r.db.Delete(&model.Product{}, id)
	if res.Error != nil {
		return fmt.Errorf("delete product: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("delete product: %w", util.ErrNotFound)
	}
	return nil
}
