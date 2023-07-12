package main

import "gin-gonic/gin"

func setupRoutes(r *gin.Engine) {
	// ruleid: route-handlers
	r.POST("/blah", blahIndexGet)

	// ruleid: route-handlers
	r.GET("/blah", blahIndexPost)

	// ruleid: route-handlers
	r.GET("/blah", blahIndexHandler)

	// ok: route-handlers
	r.POST("/blah", blahIndexPost)

	// ok: route-handlers
	r.GET("/blah", blahIndexGet)
}
