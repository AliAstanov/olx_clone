package handlers

import (
	"log"
	"strconv"
	"time"

	helpers "github.com/AliAstanov/helper"
	"github.com/AliAstanov/olx_clone/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) CreateAdmin(ctx *gin.Context) {
	var reqBody models.CreateAdmin
	var admin = &models.Admins{}

		// Request body ni JSON formatida olish
	if err := ctx.BindJSON(&reqBody); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body on create admin"})
		log.Println("Invalid request body on create admin:", err)
		return
	}

	if err := helpers.DataParser1(reqBody, admin); err != nil {
		ctx.JSON(500, gin.H{"error": "invalid request body"})
		log.Println("Error parsing request body:", err)
		return
	}

	admin.Id = uuid.New()
	admin.CreatedAt = time.Now()
	admin.UpdatedAt = time.Now()

	_, err := h.storage.GetAdminsRepo().CreateAdmin(ctx, admin)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to create admin"})
		log.Println("Failed to create admin:", err)
		return
	}

	ctx.JSON(201, gin.H{"message": "Admin created successfully"})
	log.Println("Admin created successfully with ID:", admin.Id)
}

func (h *handler) GetAdmins(ctx *gin.Context) {
	
	var reqBody models.GetListReq
	var err error

	limit := ctx.Query("limit")
	page := ctx.Query("page")

	// "limit" ni integerga aylantirish
	reqBody.Limit, err = strconv.Atoi(limit)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid 'limit' parameter value"})
		log.Println("Invalid 'limit' parameter:", err)
		return
	}
	// "page" ni integerga aylantirish
	reqBody.Page, err = strconv.Atoi(page)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid 'page' parameter value"})
		log.Println("Invalid 'page' parameter:", err)
		return
	}

	admins, err := h.storage.GetAdminsRepo().GetAdmins(ctx,&reqBody)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to fetch admins list"})
		log.Println("Failed to fetch admins list:", err)
		return
	}

	ctx.JSON(200,admins)
}

func(h *handler)GetAdminById(ctx *gin.Context){
	id := ctx.Param("id")

	admin, err := h.storage.GetAdminsRepo().GetAdminById(ctx,id)
	if err != nil {
		ctx.JSON(500,gin.H{"error":"Failed get admin by id"})
		log.Println("Failed to get admin by id:",err)
		return 
	}

	ctx.JSON(200,admin)
}
func(h *handler)UpdateAdmin(ctx *gin.Context){
	var reqBody models.UpdateAdmin
	id := ctx.Param("id")

	// Request body ni JSON formatida olish
	if err := ctx.BindJSON(&reqBody); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
		log.Println("Invalid request body:", err)
		return
	}


	// Adminni ma'lumotlar bazasida yangilash
	admin, err := h.storage.GetAdminsRepo().UpdateAdmin(ctx, &reqBody, id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to update admin"})
		log.Println("Failed to update admin:", err)
		return
	}

	ctx.JSON(200,admin)
}

func(h *handler)DeleteAdmin(ctx *gin.Context){
	id := ctx.Param("id")

	err := h.storage.GetAdminsRepo().DeleteAdmin(ctx,id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to delete admin"})
		log.Println("Failed to delete admin:", err)
		return
	
	}

	ctx.JSON(200, gin.H{"message": "Admin deleted successfully"})

}


