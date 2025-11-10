package postgres

import (
	"context"
	"log"

	helpers "github.com/AliAstanov/helper"
	"github.com/AliAstanov/olx_clone/models"
	repoi "github.com/AliAstanov/olx_clone/storage/repoI"
	"github.com/jackc/pgx/v5"
)

type UserRepo struct {
	db *pgx.Conn
}

func NewUserRepo(db *pgx.Conn) repoi.UserRepoI {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) CreateUsers(ctx context.Context, req *models.User) (*models.User, error) {

	query := `
		INSERT INTO 
			users(
				id,
				name,
				email,
				password,
				created_at,
				deleted_at
			)VALUES(
				$1,$2,$3,$4,$5,$6
			)
	`
	_, err := u.db.Exec(ctx, query,
		req.UserId,
		req.Name,
		req.Email,
		req.Password,
		req.CreatedAt,
		req.DeletedAt,
	)
	if err != nil {
		log.Println("error on creating user:", err)
		return nil, err
	}

	userIdStr := req.UserId.String() // uuid to string
	user, err := u.GetUserByid(ctx, userIdStr)

	if err != nil {
		log.Println("fail return user:", err)
		return nil, err
	}

	return user, nil
}

func (u *UserRepo) GetUsers(ctx context.Context, req *models.GetListReq) (*models.GetListUsers, error) {

	limit := req.Limit
	if limit == 0 {
		limit = DefaultLimit
	}
	page := req.Page
	if page == 0 {
		page = DefaultPage
	}

	offset := helpers.Offset(limit, page)
	// offset := (page - 1) * limit

	query := `
		SELECT 
			id,
			name,
			email,
			password,
			created_at
		FROM
			users
		WHERE 
			deleted_at IS NULL
		LIMIT 
			$1
		OFFSET 
			$2
	`
	rows, err := u.db.Query(ctx, query, limit, offset)
	if err != nil {
		log.Println("eroor on GetListUsers:", err)
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.UserId,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
		); err != nil {
			log.Println("error on scanning users row:", err)
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		log.Println("error on rows in GetListUsers:", err)
		return nil, err
	}

	return &models.GetListUsers{
		Users: users,
		Count: len(users),
	}, nil
}

func (u *UserRepo) GetUserByid(ctx context.Context, id string) (*models.User, error) {
	var user models.User

	query := `
		SELECT 
			*
		FROM
			users
		WHERE
			id = $1
	`
	err := u.db.QueryRow(
		ctx, query, id,
	).Scan(
		&user.UserId,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		log.Println("error on GetUser:", err)
		return nil, err
	}
	return &user, nil
}
func (u *UserRepo) UpdateUser(ctx context.Context, req *models.UpdateUserReq, id string) (*models.User, error) {

	query := `
        UPDATE
    		users
		SET
    		name = COALESCE($1, name),
    		email = COALESCE($2, email),
    		password = COALESCE($3, password)
		WHERE
    		id = $4;
    `
	_, err := u.db.Exec(ctx, query, req.Name, req.Email, req.Password, id)
	if err != nil {
		log.Printf("Error updating user with ID %s: %v", id, err)
		return nil, err
	}

	user, err := u.GetUserByid(ctx, id)
	if err != nil {
		log.Printf("Error fetching updated user with ID %s: %v", id, err)
		return nil, err
	}

	return user, nil
}

func (u *UserRepo) DeleteUser(ctx context.Context, id string) error {
	query := `
		UPDATE 
			users
		SET 
			deleted_at = CURRENT_TIMESTAMP
		WHERE 
			id = $1
	`
	_, err := u.db.Exec(ctx, query, id)
	if err != nil {
		log.Println("Error on soft delete:", err)
		return err
	}

	return nil
}

func (u *UserRepo) ArchiveDeleteUsers(ctx context.Context) error {
	// Tranzaksiyani boshlash
	tx, err := u.db.Begin(ctx)
	if err != nil {
		log.Println("Failed to start transaction:", err)
		return err
	}
	defer tx.Rollback(ctx) // Xato ro‘y berganda tranzaksiyani bekor qilish

	// 1. Arxivga nusxa ko'chirish
	copyQuery := `
		INSERT INTO 
			users_archive (
				id, 
				name, 
				email, 
				password, 
				created_at, 
				deleted_at
			)SELECT 
				id, 
				name, 
				email,
				password, 
				created_at, 
				deleted_at
			FROM 
				users
			WHERE 
				deleted_at IS NOT NULL
	`
	_, err = tx.Exec(ctx, copyQuery)
	if err != nil {
		log.Println("Error copying to archive:", err)
		return err
	}

	// 2. Asosiy jadvaldan o'chirish
	deleteQuery := `
		DELETE FROM 
			users
		WHERE 
			deleted_at IS NOT NULL
	`
	_, err = tx.Exec(ctx, deleteQuery)
	if err != nil {
		log.Println("Error deleting from main table:", err)
		return err
	}

	// Tranzaksiyani saqlash
	if err := tx.Commit(ctx); err != nil {
		log.Println("Transaction commit failed:", err)
		return err
	}

	log.Println("Archived and deleted users successfully.")
	return nil
}

//                bu archive va delete ning oddiy usuli
// func (u *UserRepo) ArchiveDeleteUsers(ctx context.Context) error {
// 	copyQuery := `
// 		INSET INTO
// 			users_archive(
// 				id,
// 				name,
// 				email,
// 				password,
// 				created_at,
// 				deleted_at
// 			)SELECT
// 				id,
// 				name,
// 				email,
// 				password,
// 				created_at,
// 				deleted_at
// 			FROM
// 				users
// 			WHERE
// 				deleted_at IS NOT NULL
// 			`
// 	_, err := u.db.Exec(ctx,copyQuery)
// 	if err != nil {
// 		log.Println("error on archive delete users:",err)
// 		return err
// 	}

// 	deleteQuery := `
// 		DELETE FROM
// 			users
// 		WHERE
// 		 	deleted_at IS NOT NULL
// 	`
// 	_, err = u.db.Exec(ctx,deleteQuery)
// 	if err != nil{
// 		log.Println("error from main table:",err)
// 		return nil
// 	}

// 	return nil
// }
