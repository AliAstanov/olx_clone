package handlers

import (
	"log"
	"strconv"
	"time"

	"github.com/AliAstanov/olx_clone/models"
	"github.com/AliAstanov/olx_clone/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) CreateProductWithImages(ctx *gin.Context) {
	// Multipart form olish (product data + files)
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid form data"})
		return
	}

	// JSON field’dan product ma’lumotini olish
	productData := form.Value["product"] // "product" field JSON string bo‘ladi
	if len(productData) == 0 {
		ctx.JSON(400, gin.H{"error": "Product data is required"})
		return
	}

	var reqBody models.CreateProductReq
	if err := utils.ParseJSON(productData[0], &reqBody); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid product JSON"})
		return
	}

	// Product struct yaratish
	product := &models.Product{}
	if err := utils.DataParser1(reqBody, product); err != nil {
		ctx.JSON(400, gin.H{"error": "Failed parsing product data"})
		return
	}

	product.ID = uuid.New()
	product.CreatedAt = time.Now()
	product.Status = models.StatusPending

	// DB ga saqlash
	createdProduct, err := h.storage.GetProductsRepo().CreateProducts(ctx, product)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to create product"})
		return
	}

	// Fayllarni olish (field name: "images")
	files := form.File["images"]
	var imageURLs []string
	for _, file := range files {
		filename := uuid.New().String() + "_" + file.Filename
		path := "assets/uploads/products/" + filename

		if err := ctx.SaveUploadedFile(file, path); err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to save image"})
			return
		}

		imageURL := "/assets/uploads/products/" + filename
		imageURLs = append(imageURLs, imageURL)

		// DB ga yozish
		_, err := h.storage.GetProductsRepo().AddProductImage(ctx, createdProduct.ID.String(), imageURL)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to save image in DB"})
			return
		}
	}

	ctx.JSON(201, gin.H{
		"message": "Product created successfully with images",
		"product": createdProduct,
		"images":  imageURLs,
	})
}

func (h *handler) GetProducts(ctx *gin.Context) {
	var reqBody models.GetListReq
	var err error

	limit := ctx.Query("limit")
	page := ctx.Query("page")

	reqBody.Limit, err = strconv.Atoi(limit)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid limit parameter"})
		log.Println("Invalid limit parameter:", err)
		return
	}
	reqBody.Page, err = strconv.Atoi(page)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid page parameter"})
		log.Println("Invalid page parameter:", err)
		return
	}

	products, err := h.storage.GetProductsRepo().GetListProducts(ctx, &reqBody)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed get products"})
		log.Println("Failed get products:", err)
		return
	}
	ctx.JSON(200, products)

}

func (h *handler) GetProductById(ctx *gin.Context) {
	id := ctx.Param("id")

	product, err := h.storage.GetProductsRepo().GetProductsById(ctx, id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to GetProductById"})
		log.Println("Failed to GetProductById:", err)
		return
	}
	ctx.JSON(200, product)

}

func (h *handler) UpdateProducts(ctx *gin.Context) {
	var reqBody models.UpdateProductReq
	id := ctx.Param("id")

	if err := ctx.BindJSON(&reqBody); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to request body in updateProduct"})
		log.Println("Failed to request body in updateProduct")
		return
	}

	UpdatedProduct, err := h.storage.GetProductsRepo().UpdateProducts(ctx, &reqBody, id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to Update Product"})
		log.Println("Failed UpdateProduct:", err)
		return
	}

	ctx.JSON(200, UpdatedProduct)
}

func (h *handler) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.storage.GetProductsRepo().DeleteProducts(ctx, id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to delete product"})
		log.Println("Failed to deleteProduct:", err)
		return
	}
}

func (h *handler) ApproveProduct(ctx *gin.Context) {

	id := ctx.Param("id")
	status := ctx.Param("status")

	productForUpdate, err := h.storage.GetProductsRepo().GetProductsById(ctx, id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Product not found"})
		log.Println("Product not found:", err)
		return
	}

	if productForUpdate.Status != models.StatusPending {
		ctx.JSON(400, gin.H{"error": "Product is not pending approval"})
		log.Println("Product is not pending approval")
		return
	}

	if status == "approved" {
		// Mahsulot holatini active ga o'zgartirish va amal qilish muddatini 1 oyga qo‘shish
		productForUpdate.Status = models.StatusActive
		idStr := productForUpdate.ID.String()

		err = h.storage.GetProductsRepo().SetStatus(ctx, idStr, models.StatusActive)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to approve product"})
			log.Println("Failed to approve product:", err)
			return
		}
	}

	if status == "reject" {
		// Mahsulot holatini reject ga o'zgartirish
		productForUpdate.Status = models.StatusRejected
		idStr := productForUpdate.ID.String()

		err = h.storage.GetProductsRepo().SetStatus(ctx, idStr, models.StatusRejected)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "The product could not be rejected"})
			log.Println("Failed to rejected product:", err)
			return
		}
	}

	ctx.JSON(200, gin.H{"message": "Product approved successfully", "product": productForUpdate})
	log.Println("Product approved successfully with ID:", productForUpdate.ID)
}
