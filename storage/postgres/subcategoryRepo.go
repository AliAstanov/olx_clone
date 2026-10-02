package postgres

import (
	"context"
	"log"

	halpers "github.com/AliAstanov/helper"
	"github.com/AliAstanov/olx_clone/models"
	repoi "github.com/AliAstanov/olx_clone/storage/repoI"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Subcategories struct {
	db *pgxpool.Pool
}

func NewSubcategory(db *pgxpool.Pool) repoi.SubcategoriesRepoI {
	return &Subcategories{db: db}
}

func (s *Subcategories) CreateSubcategory(ctx context.Context, req *models.Subcategories) (*models.Subcategories, error) {
	query := `
		INSERT INTO 
			subcategories(
				id,
				name,
				category_id
			)VALUES(
			$1,$2,$3
			)`
	_, err := s.db.Exec(
		ctx, query,
		req.Id,
		req.Name,
		req.CategoryID,
	)
	if err != nil {
		log.Println("error on create subcategoryRepo:", err)
		return nil, err
	}

	stringId := req.Id.String()
	subcategory, err := s.GetSubcategoriesById(ctx, stringId)
	if err != nil {
		log.Println("error GetSubCategory for CreateSubcategory:", err)
		return nil, err
	}

	return subcategory, nil
}
func (s *Subcategories) GetListSubcategories(ctx context.Context, req *models.GetListReq) (*models.GetListSubcategories, error) {
	query := `
		SELECT 
			id,
			name,
			category_id
		FROM
			subcategories
		LIMIT 
			$1
		OFFSET
			$2
	`
	limit := req.Limit
	if limit == 0 {
		limit = DefaultLimit
	}
	page := req.Page
	if page == 0 {
		page = DefaultPage
	}
	offset := halpers.Offset(limit, page)

	rows, err := s.db.Query(ctx, query, limit, offset)
	if err != nil {
		log.Println("error on GEtListSubcategories:", err)
		return nil, err
	}
	defer rows.Close()

	var subcategories []models.Subcategories
	for rows.Next() {
		var subcategory models.Subcategories
		if err := rows.Scan(
			&subcategory.Id,
			&subcategory.Name,
			&subcategory.CategoryID,
		); err != nil {
			log.Println("error on scanning subcategory rows:", err)
			return nil, err
		}
		subcategories = append(subcategories, subcategory)
	}

	return &models.GetListSubcategories{
		Subcategories: subcategories,
		Count:         len(subcategories),
	}, nil
}
func (s *Subcategories) GetSubcategoriesById(ctx context.Context, id string) (*models.Subcategories, error) {
	query := `
		SELECT
			id,
			name,
			category_id
		FROM 
			subcategories
		WHERE 
			id = $1
	`
	var subcategory models.Subcategories
	if err := s.db.QueryRow(
		ctx, query, id,
	).Scan(
		&subcategory.Id,
		&subcategory.Name,
		&subcategory.CategoryID,
	); err != nil {
		log.Println("error on GetSubcategoriesByid:", err)
		return nil, err
	}
	return &subcategory, nil
}
func (s *Subcategories) UpdateSubcategories(ctx context.Context, req *models.UpdateSubcategories, id string) (*models.Subcategories, error) {
	query := `
		UPDATE
			subcategories
		SET
			name = $1,
			category_id = $2
		WHERE
			id = $3
	`
	_, err := s.db.Exec(ctx, query, id)
	if err != nil {
		log.Println("error on update subcategories:", err)
		return nil, err
	}
	subcategory, err := s.GetSubcategoriesById(ctx, id)
	if err != nil {
		log.Println("error on GetSubcategoriesById for updating:", err)
		return nil, err
	}

	return subcategory, nil
}
func (s *Subcategories) DeleteSubcategories(ctx context.Context, id string) error {
	query := `
		DELETE FROM
			subcategories
		WHERE
			id = $1
	`
	_, err := s.db.Exec(ctx, query, id)
	if err != nil {
		log.Println("erroro n delete subcategory:", err)
		return err
	}
	return nil
}
