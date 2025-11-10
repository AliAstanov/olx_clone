package repoi

import (
	"context"

	"github.com/AliAstanov/olx_clone/models"
)

type SubcategoriesRepoI interface {
	CreateSubcategory(ctx context.Context, req *models.Subcategories) (*models.Subcategories, error)
	GetListSubcategories(ctx context.Context, req *models.GetListReq) (*models.GetListSubcategories, error)
	GetSubcategoriesById(ctx context.Context, id string) (*models.Subcategories, error)
	UpdateSubcategories(ctx context.Context, req *models.UpdateSubcategories, id string) (*models.Subcategories, error)
	DeleteSubcategories(ctx context.Context, id string) error
}
