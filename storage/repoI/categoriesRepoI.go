package repoi

import (
	"context"

	"github.com/AliAstanov/olx_clone/models"
)

type CategoriesRepoI interface {
	CreateCategory(ctx context.Context, req *models.Categories) (*models.Categories, error)
	GetListCategories(ctx context.Context, req *models.GetListReq) (*models.GetListCategories, error) 
	GetCategoriesById(ctx context.Context, id string) (*models.Categories, error)
	UpdateCategories(ctx context.Context, req *models.UpdateCategories, id string) (*models.Categories, error)
	DeleteCategories(ctx context.Context, id string) error
}
