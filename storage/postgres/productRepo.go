package postgres

import (
	"context"
	"errors"
	"log"
	"time"

	halpers "github.com/AliAstanov/helper"
	"github.com/AliAstanov/olx_clone/models"
	repoi "github.com/AliAstanov/olx_clone/storage/repoI"
	"github.com/jackc/pgx/v5"
)

type ProductRepo struct {
	db *pgx.Conn
}

func NewProductRepo(db *pgx.Conn) repoi.ProductRepoI {
	return &ProductRepo{db: db}
}

func (p *ProductRepo) CreateProducts(ctx context.Context, req *models.Product) (*models.Product, error) {
	query := `
		INSERT INTO 
			products(
				id,
				user_id,
				title,
				description,
				price,
				status,
				created_at,
				subcategory_id,
				image,
				is_new,
				expires_at
		)VALUES(
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11
		)`
	_, err := p.db.Exec(
		ctx, query,
		req.ID,
		req.UserID,
		req.Title,
		req.Description,
		req.Price,
		req.Status,
		req.CreatedAt,
		req.SubcategoryID,
		req.Image,
		req.IsNew,
		req.ExpiresAt,
	)
	if err != nil {
		log.Println("error on createProducts:", err)
		return nil, err
	}

	strId := req.ID.String()
	product, err := p.GetProductsById(ctx, strId)
	if err != nil {
		log.Println("error on CreateProduct: GetProduct for create Method:", err)
		return nil, err
	}

	return product, nil
}
func (p *ProductRepo) GetListProducts(ctx context.Context, req *models.GetListReq) (*models.GetListProducts, error) {
	limit := req.Limit
	if limit == 0 {
		limit = DefaultLimit
	}
	page := req.Page
	if page == 0 {
		page = DefaultPage
	}

	offset := int(halpers.Offset(limit, page))

	query := `
		SELECT 
			id,
			user_id,
			title,
			description,	
			price,
			status,
			created_at,		
			subcategory_id,
			image,
			is_new,
			expires_at
		FROM
			products
		OFFSET
			$1
		limit
			$2
		)`
	rows, err := p.db.Query(ctx, query, limit, offset)
	if err != nil {
		log.Println("error on GetListProducts:", err)
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(
			&product.ID,
			&product.UserID,
			&product.Title,
			&product.Description,
			&product.Price,
			&product.Status,
			&product.CreatedAt,
			&product.SubcategoryID,
			&product.Image,
			&product.IsNew,
			&product.ExpiresAt,
		); err != nil {
			log.Println("erroron GetListProducts:", err)
			return nil, err
		}
		products = append(products, product)
	}

	return &models.GetListProducts{
		Products: products,
		Count:    len(products),
	}, nil
}
func (p *ProductRepo) GetProductsById(ctx context.Context, id string) (*models.Product, error) {
	query := `
		SELECT
			id,
			user_id,
			title,
			description,	
			price,
			status,
			created_at,		
			subcategory_id,
			image,
			is_new,
			expires_at
		FROM
			products
		WHERE
			id = $1
	`
	var product models.Product
	err := p.db.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.UserID,
		&product.Title,
		&product.Description,
		&product.Price,
		&product.Status,
		&product.CreatedAt,
		&product.SubcategoryID,
		&product.Image,
		&product.IsNew,
		&product.ExpiresAt,
	)
	if err != nil {
		log.Println("error on GetProductById:", err)
		return nil, err
	}
	return &product, nil
}
func (p *ProductRepo) UpdateProducts(ctx context.Context, req *models.UpdateProductReq, id string) (*models.Product, error) {
	query := `
		UPDATE
			products
		SET 
			title = $1,
			description = $2,
			price = $3,
			image = $4,
			is_new = $5,
			expires_at = $6
		WHERE
			id = $7
	`
	// Old productni olish
	oldProduct, err := p.GetProductsById(ctx, id)
	if err != nil {
		log.Println("error getting old product:", err)
		return nil, err
	}

	// Eski qiymatlarni to'ldirish
	if req.Title == "" {
		req.Title = oldProduct.Title
	}
	if req.Description == "" {
		req.Description = oldProduct.Description
	}
	if req.Price == 0.0 {
		req.Price = oldProduct.Price
	}
	if req.Image == "" {
		req.Image = oldProduct.Image
	}
	req.IsNew = req.IsNew || oldProduct.IsNew
	if req.ExpiresAt.IsZero() {
		req.ExpiresAt = oldProduct.ExpiresAt
	}

	// Yangilash
	_, err = p.db.Exec(ctx, query,
		req.Title,
		req.Description,
		req.Price,
		req.Image,
		req.IsNew,
		req.ExpiresAt,
		id, // ID oxirgi parametr
	)
	if err != nil {
		log.Println("error on UpdateProducts:", err)
		return nil, err
	}

	// Yangilangan productni olish
	UpdatedProduct, err := p.GetProductsById(ctx, id)
	if err != nil {
		log.Println("error getting updated product:", err)
		return nil, err
	}

	return UpdatedProduct, nil
}

func (p *ProductRepo) DeleteProducts(ctx context.Context, id string) error {
	query := `
		DELETE FROM 
			products
		WHERE
			id = $1
	`
	_, err := p.db.Exec(ctx, query, id)
	if err != nil {
		log.Println("error on DeleteProduct:", err)
		return err
	}

	return nil
}

func (p *ProductRepo) SetStatus(ctx context.Context, productID string, status models.ProductStatus) error {

	switch status {
	case models.StatusPending, models.StatusActive, models.StatusRejected, models.StatusExpires:
		
		var expiresAt *time.Time
		
		if status == models.StatusActive {
			tempExpiresAt := time.Now().AddDate(0, 1, 0)
			expiresAt = &tempExpiresAt
		}

		if status == models.StatusRejected {
			// Bu holatda hech qanday amal qilish muddati qo'shilmaydi
			expiresAt = nil
		}
		// Bazani yangilash
		query := `
			UPDATE
				products
			SET
				status = $1,
				expires_at = $2
			WHERE
				id = $3
		`
		_, err := p.db.Exec(ctx, query, status, expiresAt, productID)
		if err != nil {
			log.Println("error on SetStatus:", err)
			return err
		}
		return nil
	default:
		// Noto'g'ri status
		return errors.New("invalid status")
	}
}
