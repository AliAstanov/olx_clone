package storage

import (
	"github.com/AliAstanov/olx_clone/storage/postgres"
	repoi "github.com/AliAstanov/olx_clone/storage/repoI"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StorageI interface {
	GetUserRepo() repoi.UserRepoI
	GetCategoriesRepo() repoi.CategoriesRepoI
	GetSubcategoriesRepo() repoi.SubcategoriesRepoI
	GetProductsRepo() repoi.ProductRepoI
	GetAdminsRepo() repoi.AdminsRepoI // Admins uchun repository
}

type storage struct {
	userRepo          repoi.UserRepoI
	categoriesRepo    repoi.CategoriesRepoI
	subcategoriesRepo repoi.SubcategoriesRepoI
	productRepo       repoi.ProductRepoI
	adminsRepo        repoi.AdminsRepoI // Admins uchun repository
}

func NewStorage(db *pgxpool.Pool) StorageI {
	return &storage{
		userRepo:          postgres.NewUserRepo(db),
		categoriesRepo:    postgres.NewCategoryRepo(db),
		subcategoriesRepo: postgres.NewSubcategory(db),
		productRepo:       postgres.NewProductRepo(db),
		adminsRepo:        postgres.NewAdminsRepo(db), // Admins uchun repo
	}
}

func (s *storage) GetUserRepo() repoi.UserRepoI {
	return s.userRepo
}

func (s *storage) GetCategoriesRepo() repoi.CategoriesRepoI {
	return s.categoriesRepo
}

func (s *storage) GetSubcategoriesRepo() repoi.SubcategoriesRepoI {
	return s.subcategoriesRepo
}

func (s *storage) GetProductsRepo() repoi.ProductRepoI {
	return s.productRepo
}

func (s *storage) GetAdminsRepo() repoi.AdminsRepoI {
	return s.adminsRepo
}
