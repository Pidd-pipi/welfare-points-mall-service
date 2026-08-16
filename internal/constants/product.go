package constants

// ProductCategory 商品类型枚举，前后端共享定义（frontend/src/constants/product.ts 对应实现）。
type ProductCategory string

const (
	CategoryPhysical ProductCategory = "physical"
	CategoryVirtual  ProductCategory = "virtual"
	CategoryService  ProductCategory = "service"
)

func (c ProductCategory) Valid() bool {
	switch c {
	case CategoryPhysical, CategoryVirtual, CategoryService:
		return true
	}
	return false
}
