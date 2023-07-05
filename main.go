package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Book struct {
	ID     uint
	Title  string
	Author string
}

// Set default Handler
// func defaultHandler(c *gin.Context) {
// 	c.HTML(http.StatusOK, "default.html", gin.H{})
// }

// Setup router
func setupRouter(r *gin.Engine, db *gorm.DB) {
	r.LoadHTMLGlob("templates/**/*.html")
	r.Use(connectDatabase(db))
	r.GET("/books/", bookIndexHandler)
	r.GET("/books/new", bookNewGetHandler)
	r.POST("/books/new", bookNewPostHandler)
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/books/")
	})
}

// Middleware for connecting to database
func connectDatabase(db *gorm.DB) gin.HandlerFunc  {
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
	err = router.Run(":80")
	if err != nil {
		log.Fatalf("gin Run error: %s", err)
	}
}

func bookIndexHandler(ctx *gin.Context)  {
	db := ctx.Value("database").(*gorm.DB)
	books := []Book{}

	if err := db.Find(&books).Error; err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.HTML(http.StatusOK, "books/index.html", gin.H{"books" : books})
}

func bookNewGetHandler(ctx *gin.Context)  {
	ctx.HTML(http.StatusOK, "books/new.html", gin.H{})
}

func bookNewPostHandler(ctx *gin.Context)  {

}