package postgres

import (
	"context"
	"log"

	helpers "github.com/AliAstanov/helper"
	"github.com/AliAstanov/olx_clone/models"
	repoi "github.com/AliAstanov/olx_clone/storage/repoI"
	"github.com/jackc/pgx/v5"
)

type AdminsRepo struct {
	db *pgx.Conn
}

func NewAdminsRepo(db *pgx.Conn) repoi.AdminsRepoI {
	return &AdminsRepo{
		db: db,
	}
}

func (a *AdminsRepo) CreateAdmin(ctx context.Context, req *models.Admins) (*models.Admins, error) {
	log.Println("Starting CreateAdmin for username:", req.Name)

	query := `
		INSERT INTO 
			admins(
				id,
				name,
				email,
				password,
				created_at,
				updated_at
			)VALUES(
				$1,$2,$3,$4,$5,$6
			)`
	_, err := a.db.Exec(
		ctx, query,
		req.Id,
		req.Name,
		req.Email,
		req.Password,
		req.CreatedAt,
		req.UpdatedAt,
	)
	if err != nil {
		log.Println("Error on CreateAdmin - failed to insert admin:", err)
		return nil, err
	}

	strId := req.Id.String()
	log.Println("CreateAdmin - successfully inserted admin with ID:", strId)

	admin, err := a.GetAdminById(ctx, strId)
	if err != nil {
		log.Println("Error on CreateAdmin - failed to retrieve admin after insert:", err)
		return nil, err
	}

	log.Println("CreateAdmin - successfully retrieved admin:", admin)
	return admin, nil
}
func (a *AdminsRepo) GetAdmins(ctx context.Context, req *models.GetListReq) (*models.GetListAdmins, error) {
	log.Println("Started GetAdmins with limit:", req.Limit, "and page:", req.Page)
	
	var admins []models.Admins

	limit := req.Limit
	if limit == 0 {
		limit = DefaultLimit
	}
	page := req.Page
	if page == 0 {
		page = DefaultPage
	}
	offset := helpers.Offset(limit, page)

	query := `
		SELECT 
			id,
			name,
			email,
			password,
			created_at,	
			updated_at
		FROM
			admins
		LIMIT
			$1
		OFFSET
			$2
	`

	rows, err := a.db.Query(ctx, query,limit,offset)
	if err != nil{
		log.Println("Error executing GetAdmins query:", err)
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var admin models.Admins
		if err := rows.Scan(
			&admin.Id,
			&admin.Name,
			&admin.Email,
			&admin.Password,
			&admin.CreatedAt,
			&admin.UpdatedAt,
		); err != nil {
			log.Println("Error scanning row in GetAdmins:", err)
			return nil, err
		}
		admins = append(admins, admin)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error in rows iteration for GetAdmins:", err)
		return nil, err
	}

	log.Println("Completed GetAdmins successfully, total admins:", len(admins))

	return &models.GetListAdmins{
		Admins: admins,
		Count: len(admins),
	}, nil
}
func (a *AdminsRepo) GetAdminById(ctx context.Context, id string) (*models.Admins, error) {
	log.Println("Started GetAdminById for ID:", id)

	query := `
		SELECT 
			id,
			name,
			email,
			password,
			created_at,
			updated_at
		FROM
			admins
		WHERE
			id = $1`
	var admin models.Admins

	err := a.db.QueryRow(ctx, query, id).Scan(
		&admin.Id,
		&admin.Name,
		&admin.Email,
		&admin.Password,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)
	if err != nil {
		log.Println("Error in GetAdminById - scanning row:", err)
		return nil, err
	}

	log.Println("Successfully fetched admin with ID:", id)
	return &admin, nil
}

func (a *AdminsRepo) UpdateAdmin(ctx context.Context, req *models.UpdateAdmin, id string) (*models.Admins, error) {
	log.Println("Started UpdateAdmin for ID:", id)

	query := `
		UPDATE
			admins
		SET
			name = $1,
			email = $2,
			password = $3
		WHERE 
			id = $4`

	_, err := a.db.Exec(ctx, query, &req.Name, &req.Email, &req.Password, id)
	if err != nil {
		log.Println("Error in UpdateAdmin - executing update:", err)
		return nil, err
	}

	admin, err := a.GetAdminById(ctx, id)
	if err != nil {
		log.Println("Error in UpdateAdmin - fetching updated admin:", err)
		return nil, err
	}

	log.Println("Successfully updated admin with ID:", id)
	return admin, nil
}

func (a *AdminsRepo) DeleteAdmin(ctx context.Context, id string) error {
	log.Println("Started DeleteAdmin for ID:", id)

	query := `
		DELETE FROM 
			admins
		WHERE
			id = $1
	`
	_, err := a.db.Exec(ctx, query, id)
	if err != nil {
		log.Println("Error in DeleteAdmin - executing delete:", err)
		return err
	}

	log.Println("Successfully deleted admin with ID:", id)
	return nil
}
