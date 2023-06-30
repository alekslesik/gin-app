package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func defaultHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "default.html", gin.H{})
}

func setupRouter(router *gin.Engine) {
	router.LoadHTMLGlob("templates/**/*.html")
	router.GET("/", defaultHandler)
}

func main() {
	router := gin.Default()

	setupRouter(router)

	err := router.Run(":80")
	if err != nil {
		log.Fatalf("gin Run error: %s", err)
	}
}

type Book struct {
	ID     uint
	Title  string
	Author string
}

func setupDatabase(db *gorm.DB) error {
	err := db.AutoMigrate(Book{})
	if err != nil {
		return fmt.Errorf("error migrating database: %s", err)
	}

	return nil
}
