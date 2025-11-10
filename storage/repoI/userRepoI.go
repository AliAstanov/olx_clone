package repoi

import (
	"context"

	"github.com/AliAstanov/olx_clone/models"
)

type UserRepoI interface {
	CreateUsers(ctx context.Context, req *models.User) (*models.User, error)
	GetUsers(ctx context.Context, req *models.GetListReq) (*models.GetListUsers, error)
	GetUserByid(ctx context.Context, id string) (*models.User, error)
	UpdateUser(ctx context.Context, req *models.UpdateUserReq,id string ) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
	ArchiveDeleteUsers(ctx context.Context)error
}
