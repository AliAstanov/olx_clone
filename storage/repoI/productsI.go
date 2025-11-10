package repoi

import (
	"context"

	"github.com/AliAstanov/olx_clone/models"
)

type ProductRepoI interface {
	CreateProducts(ctx context.Context, req *models.Product) (*models.Product, error)
	GetListProducts(ctx context.Context, req *models.GetListReq) (*models.GetListProducts, error)
	GetProductsById(ctx context.Context, id string) (*models.Product, error)
	UpdateProducts(ctx context.Context, req *models.UpdateProductReq, id string) (*models.Product, error)
	DeleteProducts(ctx context.Context, id string) error
	SetStatus(ctx context.Context, productID string, status models.ProductStatus) error
}
