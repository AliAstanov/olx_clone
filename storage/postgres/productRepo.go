package postgres

import (
	"context"
	"errors"
	"log"
	"time"

	halpers "github.com/AliAstanov/helper"
	"github.com/AliAstanov/olx_clone/models"
	repoi "github.com/AliAstanov/olx_clone/storage/repoI"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepo struct {
	db *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) repoi.ProductRepoI {
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
func (p *ProductRepo) GetProductsById(ctx context.Context, id string) (*models.Product, error) {
	var product models.Product
	var expiresAt *time.Time

	query := `
        SELECT id, user_id, title, description, price, status, created_at, subcategory_id, image, is_new, expires_at
        FROM products
        WHERE id = $1
    `
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
		&expiresAt,
	)
	if err != nil {
		log.Println("error on GetProductById:", err)
		return nil, err
	}

	// ExpiresAt pointerga tayinlash (NULL bo'lsa nil)
	product.ExpiresAt = expiresAt

	// Product images olish
	imgQuery := `SELECT image_url FROM product_images WHERE product_id = $1`
	rows, err := p.db.Query(ctx, imgQuery, id)
	if err != nil {
		log.Println("error on GetProductById images:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			log.Println("error scanning product images:", err)
			return nil, err
		}
		product.Images = append(product.Images, url)
	}

	return &product, nil
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
        SELECT id, user_id, title, description, price, status, created_at, subcategory_id, image, is_new, expires_at
        FROM products
        ORDER BY created_at DESC
        OFFSET $1 LIMIT $2
    `
	rows, err := p.db.Query(ctx, query, offset, limit)
	if err != nil {
		log.Println("error on GetListProducts:", err)
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var pdt models.Product
		var expiresAt *time.Time

		if err := rows.Scan(
			&pdt.ID,
			&pdt.UserID,
			&pdt.Title,
			&pdt.Description,
			&pdt.Price,
			&pdt.Status,
			&pdt.CreatedAt,
			&pdt.SubcategoryID,
			&pdt.Image,
			&pdt.IsNew,
			&expiresAt,
		); err != nil {
			log.Println("error scanning product:", err)
			return nil, err
		}

		pdt.ExpiresAt = expiresAt

		// Product images olish
		imgRows, err := p.db.Query(ctx, "SELECT image_url FROM product_images WHERE product_id = $1", pdt.ID)
		if err != nil {
			log.Println("error fetching images:", err)
			return nil, err
		}

		var imgs []string
		for imgRows.Next() {
			var url string
			if err := imgRows.Scan(&url); err != nil {
				log.Println("error scanning images:", err)
				imgRows.Close()
				return nil, err
			}
			imgs = append(imgs, url)
		}
		imgRows.Close()
		pdt.Images = imgs

		products = append(products, pdt)
	}

	return &models.GetListProducts{
		Products: products,
		Count:    len(products),
	}, nil
}

func (p *ProductRepo) UpdateProducts(ctx context.Context, req *models.UpdateProductReq, id string) (*models.Product, error) {
	query := `
        UPDATE products
        SET title = $1,
            description = $2,
            price = $3,
            image = $4,
            is_new = $5,
            expires_at = $6
        WHERE id = $7
    `

	// Agar foydalanuvchi ExpiresAt yuborgan bo‘lsa, pointerga aylantiramiz
	var expiresAt *time.Time
	if !req.ExpiresAt.IsZero() {
		expiresAt = &req.ExpiresAt
	}

	_, err := p.db.Exec(ctx, query,
		req.Title,
		req.Description,
		req.Price,
		req.Image,
		req.IsNew,
		expiresAt,
		id,
	)
	if err != nil {
		log.Println("error on UpdateProducts:", err)
		return nil, err
	}

	// Yangilangan productni olish
	updatedProduct, err := p.GetProductsById(ctx, id)
	if err != nil {
		log.Println("error getting updated product:", err)
		return nil, err
	}

	return updatedProduct, nil
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

func (p *ProductRepo) AddProductImage(ctx context.Context, productId, imageUrl string) (string, error) {
	var id string

	query := `
        INSERT INTO product_images (id, product_id, image_url)
        VALUES ($1, $2, $3)
        RETURNING id
    `

	err := p.db.QueryRow(ctx, query,
		uuid.New().String(),
		productId,
		imageUrl,
	).Scan(&id)

	return id, err
}
