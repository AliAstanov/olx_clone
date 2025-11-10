package handlers

import (
	"log"
	"strconv"

	helpers "github.com/AliAstanov/helper"
	"github.com/AliAstanov/olx_clone/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) CreateSubcategories(ctx *gin.Context) {
	var reqBody models.CreateSubcategories
	var subcategories models.Subcategories

	if err := ctx.BindJSON(&reqBody); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
		log.Println("Invalid request body")
		return
	}

	if err := helpers.DataParser1(&reqBody, &subcategories); err != nil {
		ctx.JSON(400, gin.H{"error": "Failed to pars data"})
		log.Println("Failed to parse request body:", err)
		return
	}
	subcategories.Id = uuid.New()

	_, err := h.storage.GetSubcategoriesRepo().CreateSubcategory(ctx, &subcategories)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed create subcategories"})
		log.Println("Failed CreateSubcategories:", err)
		return
	}
	ctx.JSON(201, gin.H{"message": "create subcategory soccessfully"})
}

func (h *handler) GetSubcategories(ctx *gin.Context) {
	var req models.GetListReq
	var err error

	limit := ctx.Query("limit")
	page := ctx.Query("page")

	req.Limit, err = strconv.Atoi(limit)
	if err != nil {
		ctx.JSON(400,gin.H{"error": "Invalid limit parameter"})
		log.Println("Invalidd limit parameter:",err)
		return
	}
	
	req.Page, err = strconv.Atoi(page)
	if err != nil {
		ctx.JSON(400,gin.H{"error":"Invalid page parameter"})
		log.Println("Invalid page parameter:",err)
		return
	}

	subcategories, err := h.storage.GetSubcategoriesRepo().GetListSubcategories(ctx,&req)
	if err != nil {
		ctx.JSON(500,gin.H{"error":"Failed get subcategories"})
		log.Println("Failed get subcategories:",err)
		return
	}
	ctx.JSON(200,subcategories)
}

func(h *handler) GetSubcategoryById(ctx *gin.Context){
	id := ctx.Param("id")

	subcategory, err :=h.storage.GetSubcategoriesRepo().GetSubcategoriesById(ctx,id)
	if err != nil {ctx.JSON(500,gin.H{"error":"Failed get subcategory"})
	log.Println("Failed get subcategory:",err)
	return
}
	ctx.JSON(200,subcategory)
}

func(h *handler) UpdateSubcategory(ctx *gin.Context){
	var reBody models.UpdateSubcategories

	if err := ctx.BindJSON(&reBody); err != nil {
		ctx.JSON(400,gin.H{"error":"Invalid bind request body"})
		log.Println("Invalid bind reqBody")
		return
	}

}

func(h *handler)DeleteSubcategory(ctx *gin.Context){
	id := ctx.Param("id")

	err := h.storage.GetSubcategoriesRepo().DeleteSubcategories(ctx,id)
	if err != nil {
		ctx.JSON(500,gin.H{"error":"Failed to delete subcategory"})
		log.Println("Failed delete subcategory",err)
		return
	}
	ctx.JSON(200,"delete subcategory soccessfully")
}