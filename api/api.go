package api

import (
	"net/http"

	"github.com/AliAstanov/olx_clone/api/handlers"
	"github.com/AliAstanov/olx_clone/storage"
	"github.com/gin-gonic/gin"
)

func Api(storage storage.StorageI) {
	router := gin.Default()

	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*")
	router.MaxMultipartMemory = 8 << 20 //8mib

	h := handlers.NewHandler(storage)

	router.GET("/ping", h.Ping)

	//engine := router.Group("/")

	//userRepo
	router.POST("/create-user", h.CreateUser)
	router.GET("/get-users", h.GetUsers)
	router.GET("/get-user/:id", h.GetUserById)
	router.PUT("/update-user/:id", h.UpdateUser)
	router.PUT("/delete-user/:id", h.DeletUser)
	router.DELETE("/archive-delete-users", h.ArchiveAndDeleteUser)

	//categoryRepo
	router.POST("/create-category", h.CreateCategories)
	router.GET("/get-categories", h.GetCategories)
	router.GET("/get-category/:id", h.GetCategoriesById)
	router.PUT("/update-category/:id", h.UpdateCategories)
	router.DELETE("/delete-category/:id", h.DeleteCategories)

	//subcategory
	router.POST("/create-subcategory", h.CreateSubcategories)
	router.GET("/get-subcategories", h.GetSubcategories)
	router.GET("/get-subcategory/:id", h.GetSubcategoryById)
	router.PUT("/update-subcategory", h.UpdateSubcategory)
	router.DELETE("/delete-subcategory/:id", h.DeleteSubcategory)

	//admin
	router.POST("/create-admin", h.CreateAdmin)
	router.GET("/get-admins", h.GetAdmins)
	router.GET("/get-admin/:id", h.GetAdminById)
	router.PUT("/update-admin/:id", h.UpdateAdmin)
	router.DELETE("/delete-admin/:id", h.DeleteAdmin)

	//product
	router.POST("/create-product", h.CreateProduct)
	router.GET("/get-products", h.GetProducts)
	router.GET("/get-product/:id", h.GetAdminById)
	router.PUT("/update-product", h.UpdateProducts)
	router.DELETE("/delete-product", h.DeleteProduct)

	//approved_and_reject
	router.PUT("/approve-product/:id/:status", h.ApproveProduct)

	//browser da pageni ochib beradi
	router.GET("/get", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"laa":"olalla"})
	})

	//post image for product
	router.POST("/", h.PostImageForProduct)

	router.Run(":8081")

}


