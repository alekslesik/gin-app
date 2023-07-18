package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Book struct {
	ID     uint   `form:"-"`
	Title  string `form:"title" binding:"required"`
	Author string `form:"author" binding:"required"`
}

// Setup router
func setupRouter(r *gin.Engine, db *gorm.DB) {
	r.LoadHTMLGlob("templates/**/*.html")
	r.Use(connectDatabase(db))
	r.GET("/books/", bookIndexGet)
	r.GET("/books/new", bookNewGet)
	r.POST("/books/new", bookNewPost)
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/books/")
	})
}

// Middleware for connecting to database
func connectDatabase(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("database", db)
	}
}

func setupDatabase(db *gorm.DB) error {
	err := db.AutoMigrate(Book{})
	if err != nil {
		return fmt.Errorf("error migrating database: %s", err)
	}

	return nil
}

func main() {
	// open database
	db, err := gorm.Open(sqlite.Open("gin-app.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}

	// setup database
	err = setupDatabase(db)
	if err != nil {
		log.Fatalf("Database setup error: %s", err)
	}

	router := gin.Default()
	setupRouter(router, db)
	err = router.Run(":3000")
	if err != nil {
		log.Fatalf("gin Run error: %s", err)
	}
}

func bookIndexGet(ctx *gin.Context) {
	db := ctx.Value("database").(*gorm.DB)
	books := []Book{}

	if err := db.Find(&books).Error; err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.HTML(http.StatusOK, "books/index.html", gin.H{"books": books})
}

func bookNewGet(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "books/new.html", gin.H{})
}

func bookNewPost(ctx *gin.Context) {
	book := &Book{}

	if err := ctx.ShouldBind(book); err != nil {
		verrs := err.(validator.ValidationErrors)
		messages := make([]string, len(verrs))
		for i, verr := range verrs {
			messages[i] = fmt.Sprintf("%s is required, but was empty.", verr.Field())
		}

		ctx.HTML(http.StatusBadRequest, "books/new.html", gin.H{"errors": messages})
		return
	}

	db := ctx.Value("database").(*gorm.DB)
	if err := db.Create(book).Error; err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusFound, "/books/")
}
