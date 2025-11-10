package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	helpers "github.com/AliAstanov/helper"
	"github.com/AliAstanov/olx_clone/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) CreateProduct(ctx *gin.Context) {
	var reqBody models.CreateProductReq
	var product = &models.Product{}

	if err := ctx.BindJSON(&reqBody); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid requst body"})
		log.Println("Invalid requestBody:", err)
		return
	}

	if err := helpers.DataParser1(reqBody, product); err != nil {
		ctx.JSON(400, gin.H{"error": "Failed Parsing data on create product"})
		log.Println("Failed Parsing data on create product:", err)
		return
	}
	product.ID = uuid.New()
	product.CreatedAt = time.Now()
	product.Status = models.StatusPending

	_, err := h.storage.GetProductsRepo().CreateProducts(ctx, product)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to Create Product"})
		log.Println("Failed to CreateProduct:", err)
		return
	}
	ctx.JSON(201, gin.H{
		"message": "Product created successfully and is awaiting admin approval.",
		"product": product,
	})
	log.Println("Product created successfully and is awaiting admin approval with ID:", product.ID)
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

func (h *handler) PostImageForProduct(ctx *gin.Context) {

	//Get th file
	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"error": "Failed to upload image",
		})
		return
	}
	// Save the file
	err = ctx.SaveUploadedFile(file, "assets/uploads/"+file.Filename)
	if err != nil {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"error": "Failed to Saved image",
		})
		return
	}

	// Render the page
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"image": "/assets/uploads/" + file.Filename,
	})

}
