package repoi

import (
	"context"

	"github.com/AliAstanov/olx_clone/models"
)

type AdminsRepoI interface {
	CreateAdmin(ctx context.Context, req *models.Admins) (*models.Admins, error)
	GetAdmins(ctx context.Context, req *models.GetListReq) (*models.GetListAdmins, error)
	GetAdminById(ctx context.Context, id string) (*models.Admins, error)
	UpdateAdmin(ctx context.Context, req *models.UpdateAdmin, id string) (*models.Admins, error)
	DeleteAdmin(ctx context.Context, id string) error
}

