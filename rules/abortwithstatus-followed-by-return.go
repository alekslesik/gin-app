package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func test1(c *gin.Context) {
	// ok: abortwithstatus-followed-by-return
	c.AbortWithError(http.StatusInternalServerError)
	return
}

func test2(c *gin.Context) {
	// ruleid: abortwithstatus-followed-by-return
	c.AbortWithError(http.StatusInternalServerError)
	log.Printf("asdf")
}

func test3(c *gin.Context) {
	if true {
		// ok: abortwithstatus-followed-by-return
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
}

func test4(c *gin.Context) {
	if false {
		// ruleid: abortwithstatus-followed-by-return
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Printf("asdf")
	}
}

func test5(c *gin.Context) {
	// ok: abortwithstatus-followed-by-return
	c.AbortWithStatusJSON(http.StatusInternalServerError)
	return
}

func test6(c *gin.Context) {
	// ruleid: abortwithstatus-followed-by-return
	c.AbortWithStatusJSON(http.StatusInternalServerError)
	log.Printf("asdf")
}

func test7(c *gin.Context) {
	// ruleid: abortwithstatus-followed-by-return
	c.AbortWithStatusJSON(http.StatusInternalServerError)
	log.Printf("other stuff in between is not allowed")
	return
}
