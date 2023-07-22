package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"math"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//go:embed templates
var tmplEmbed embed.FS

//go:embed static
var staticEmbedFS embed.FS

type staticFS struct {
	fs fs.FS
}

func (sfs *staticFS) Open(name string) (fs.File, error) {
	return sfs.fs.Open(filepath.Join("static", name))
}

var staticEmbed = &staticFS{staticEmbedFS}

type Book struct {
	ID     uint   `form:"-"`
	Title  string `form:"title" binding:"required"`
	Author string `form:"author" binding:"required"`
}

func setupRouter(r *gin.Engine, db *gorm.DB) {
	tmpl := template.Must(template.ParseFS(tmplEmbed, "templates/*/*.html"))
	r.SetHTMLTemplate(tmpl)

	r.Use(connectDatabase(db))
	r.StaticFS("/static", http.FS(staticEmbed))
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

	pageStr := ctx.DefaultQuery("page", "1")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var bookCount int64

	if err := db.Table("books").Count(&bookCount).Error; err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	const booksPerPage = 15

	pageCount := int(math.Ceil(float64(bookCount) / float64(booksPerPage)))
	if pageCount == 0 {
		pageCount = 1
	}

	if page < 1 || page > pageCount {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	offset := (page - 1) * booksPerPage

	books := []Book{}

	if err := db.Limit(booksPerPage).Offset(offset).Find(&books).Error; err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var prevPage, nextPage string

	if page > 1 {
		prevPage = fmt.Sprintf("%d", page-1)
	}

	if page < pageCount {
		nextPage = fmt.Sprintf("%d", page+1)
	}

	pages := make([]int, pageCount)

	for i := 0; i < pageCount; i++ {
		pages[i] = i + 1
	}

	ctx.HTML(http.StatusOK, "books/index.html", gin.H{
		"books":     books,
		"pageCount": pageCount,
		"page":      page,
		"prevPage":  prevPage,
		"nextPage":  nextPage,
		"pages":     pages,
	})
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
