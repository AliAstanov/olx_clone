package postgres

import (
	"context"
	"log"

	helpers "github.com/AliAstanov/helper"
	"github.com/AliAstanov/olx_clone/models"
	repoi "github.com/AliAstanov/olx_clone/storage/repoI"
	"github.com/jackc/pgx/v5"
)

type CategoriesRepo struct {
	db *pgx.Conn
}

func NewCategoryRepo(db *pgx.Conn) repoi.CategoriesRepoI {
	return &CategoriesRepo{
		db: db,
	}
}

func (c *CategoriesRepo) CreateCategory(ctx context.Context, req *models.Categories) (*models.Categories, error) {

	query := `
		INSERT INTO
			categories(
				id,
				name
			)VALUES(
				$1,$2
			)
	`
	_, err := c.db.Exec(ctx, query, req.Id, req.Name)
	if err != nil {
		log.Println("error on create categories:", err)
		return nil, err
	}

	categoryIdStr := req.Id.String()
	category, err := c.GetCategoriesById(ctx, categoryIdStr)
	if err != nil {
		log.Println("fail return user:", err)
		return nil, err
	}

	return category, nil
}

func (c *CategoriesRepo) GetListCategories(ctx context.Context, req *models.GetListReq) (*models.GetListCategories, error) {
	var categories []*models.Categories

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
			name
		FROM 
			categories
		LIMIT
			$1
		OFFSET
			$2
	`
	rows, err := c.db.Query(ctx, query, limit, offset)
	if err != nil {
		log.Println("error on GetCategoriesList:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category models.Categories
		if err := rows.Scan(
			&category.Id,
			&category.Name,
		); err != nil {
			log.Println("Error in scanning categories:", err)
			return nil, err
		}
		categories = append(categories, &category)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error on rows in GetStudentList:", err)
		return nil, err
	}
	return &models.GetListCategories{
		Categories: categories,
		Count:      len(categories),
	}, nil

}

func (c *CategoriesRepo) GetCategoriesById(ctx context.Context, id string) (*models.Categories, error) {
	query := `
		SELECT 
			id,
			name
		FROM
			categories
		WHERE
			id = $1
	`
	var categories models.Categories

	err := c.db.QueryRow(ctx, query, id).Scan(
		&categories.Id,
		&categories.Name,
	)
	if err != nil {
		log.Println("error on GetCategoriesById:", err)
		return nil, err
	}

	return &categories, nil
}

func (c *CategoriesRepo) UpdateCategories(ctx context.Context, req *models.UpdateCategories, id string) (*models.Categories, error) {
	query := `
		UPDATE
			categories
		SET
			name =  $1
		WHERE
			id = $2	
	`
	_, err := c.db.Exec(ctx, query, req.Name,id)
	if err != nil {
		log.Println("error on UpdateCategories:", err)
		return nil, err
	}
	categories, err := c.GetCategoriesById(ctx, id)
	if err != nil {
		log.Println("error on GetCategoriesById for Update categories:", err)
		return nil, err
	}
	return categories, nil
}

func (c *CategoriesRepo) DeleteCategories(ctx context.Context, id string) error {
	query := `
		DELETE FROM 
			categories
		WHERE
			id = $1
	`
	_, err := c.db.Exec(ctx, query, id)
	if err != nil {
		log.Println("error on DeleteCategories:", err)
		return err
	}
	return nil
}
