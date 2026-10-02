package handlers

import (
	"log"
	"strconv"

	"github.com/AliAstanov/olx_clone/models"
	"github.com/AliAstanov/olx_clone/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) CreateCategories(ctx *gin.Context) {
	var reqBody models.CreateCategories
	var category = &models.Categories{}

	if err := ctx.BindJSON(&reqBody); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid Requestbody on create categories"})
		log.Println("Invalid request body on create categories:", err)
		return
	}

	err := utils.DataParser1(reqBody, category)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "invalid request body"})
		log.Println("Invalid request body:", err)
		return
	}

	category.Id = uuid.New()

	_, err = h.storage.GetCategoriesRepo().CreateCategory(ctx, category)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to create categories"})
		log.Println("Failed to create categories:", err)
		return
	}
	ctx.JSON(201, gin.H{"message": "Create category soccessfully"})

}

func (h *handler) GetCategories(ctx *gin.Context) {
	var reqBody models.GetListReq
	var err error

	limit := ctx.Query("limit")
	page := ctx.Query("page")

	reqBody.Limit, err = strconv.Atoi(limit)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body limit"})
		log.Println("Invalid request body limit:", err)
		return
	}
	reqBody.Page, err = strconv.Atoi(page)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body page"})
		log.Println("Invalid request body page:", err)
		return
	}

	categories, err := h.storage.GetCategoriesRepo().GetListCategories(ctx, &reqBody)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed categories list"})
		log.Println("Failed categories list:", err)
		return
	}

	ctx.JSON(200, categories)

}

func (h *handler) GetCategoriesById(ctx *gin.Context) {
	id := ctx.Param("id")

	category, err := h.storage.GetCategoriesRepo().GetCategoriesById(ctx, id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed get category by id"})
		log.Println("Failed to get category by id:", err)
		return
	}

	ctx.JSON(200, category)
}

func (h *handler) UpdateCategories(ctx *gin.Context) {
	id := ctx.Param("id")
	var reqBody models.UpdateCategories

	if err := ctx.BindJSON(&reqBody); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body to update categories"})
		log.Println("Invalid request body to update categories:", err)
		return
	}

	categories, err := h.storage.GetCategoriesRepo().UpdateCategories(ctx, &reqBody, id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed update categories"})
		log.Println("Failed to update categories:", err)
		return
	}

	ctx.JSON(200, categories)
}

func (h *handler) DeleteCategories(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.storage.GetCategoriesRepo().DeleteCategories(ctx, id)
	if err != nil {
		ctx.JSON(500, gin.H{"eerror": "Failed Delete to Category"})
		log.Println("Failed delete to Category:", err)
		return
	}
	ctx.JSON(200, gin.H{"message": "delete categories soccessfully"})
}
