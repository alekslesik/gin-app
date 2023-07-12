package main

import (
	"github.com/gin-gonic/gin"
)

// missing http method
// ruleid: handler-naming
func fooIndex(c *gin.Context) {}

// missing crudAction
// ruleid: handler-naming
func bookGet(c *gin.Context) {}

// useless repetition
// ruleid: handler-naming
func thingDeleteDelete(c *gin.Context) {}

// XXX consider allowing exported handlers
// ruleid: handler-naming
func ThingIndexGet(c *gin.Context) {}

// ok: handler-naming
func thingIndexGet(c *gin.Context) {}

// ok: handler-naming
func thingShowGet(c *gin.Context) {}

// ok: handler-naming
func thingNewPost(c *gin.Context) {}

// ok: handler-naming
func thingEditPatch(c *gin.Context) {}

// ok: handler-naming
func thingDelete(c *gin.Context) {}

// ok: handler-naming
func AnythingGoes() {}
