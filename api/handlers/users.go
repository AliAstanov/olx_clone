package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/AliAstanov/olx_clone/models"
	"github.com/AliAstanov/olx_clone/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) CreateUser(ctx *gin.Context) {

	var reqBody models.CreateUserReq
	var user = &models.User{}

	if err := ctx.BindJSON(&reqBody); err != nil {
		log.Println("Invalid request body on create users:", err)
		ctx.JSON(400, gin.H{"message": "invalid request body"})
		return
	}

	// Parolni xeshlash
	HashedPassword, err := utils.HashPassword(reqBody.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
		return
	}

	if err := utils.DataParser1(reqBody, user); err != nil {
		log.Print("Failed to parse request body:", err)
		ctx.JSON(400, gin.H{"error": "Failed to parse request body"})
		return
	}
	log.Printf("Parsed user data: %+v\n", user) // Ma'lumotlar to'g'ri kiritilganini tekshirish uchun

	user.UserId = uuid.New()
	user.CreatedAt = time.Now()
	user.Password = HashedPassword
log.Println("---===-=-=-=-=-=-=--ssss")
	_, err = h.storage.GetUserRepo().CreateUsers(ctx, user)
	if err != nil {
		log.Println("Failed to create users:", err)
		ctx.JSON(500, gin.H{"error": "Failed to create users"})
		return
	}

	ctx.JSON(201, gin.H{"message": "Created user soccessfully"})

}

func (h *handler) GetUsers(ctx *gin.Context) {
	var reqBody models.GetListReq
	var err error

	limit := ctx.Query("limit")
	page := ctx.Query("page")

	if reqBody.Limit, err = strconv.Atoi(limit); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid limit parametr"})
		log.Println("Invalid limit parametr:", err)
		return
	}

	if reqBody.Page, err = strconv.Atoi(page); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid page parametr"})
		log.Println("Invalid page parametr:", err)
	}

	users, err := h.storage.GetUserRepo().GetUsers(ctx, &reqBody)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to get user list"})
		log.Println("Failed to get user list", err)
		return
	}

	ctx.JSON(200, users)

}

func (h *handler) GetUserById(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := h.storage.GetUserRepo().GetUserByid(ctx, id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to get user by ID"})
		log.Println("Failed to get user by ID:", err)
		return
	}

	ctx.JSON(200, user)
}

func (h *handler) UpdateUser(ctx *gin.Context) {
	var reqBody models.UpdateUserReq
	id := ctx.Param("id")

	if err := ctx.BindJSON(&reqBody); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
		log.Println("Invalid request body:", err)
		return
	}

	user, err := h.storage.GetUserRepo().UpdateUser(ctx, &reqBody, id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to update user"})
		log.Println("Failed to update user:", err)
		return
	}

	ctx.JSON(200, user)
}

func (h *handler) DeletUser(ctx *gin.Context) {

	id := ctx.Param("id")

	err := h.storage.GetUserRepo().DeleteUser(ctx, id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed on delete user"})
		log.Println("Failed on delete user:", err)
		return
	}
	ctx.JSON(200, "delete user soccessfully")
}

func (h *handler) ArchiveAndDeleteUser(ctx *gin.Context) {

	err := h.storage.GetUserRepo().ArchiveDeleteUsers(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed archive user"})
		log.Println("Failed archive user")
		return
	}
	ctx.JSON(200, gin.H{"message": "Archive user soccessfully"})
}
