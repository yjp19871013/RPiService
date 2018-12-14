package router

import (
	"github.com/gin-gonic/gin"
)

var (
	getRouter = map[string][]gin.HandlerFunc{}

	postRouter = map[string][]gin.HandlerFunc{}

	patchRouter = map[string][]gin.HandlerFunc{}

	deleteRouter = map[string][]gin.HandlerFunc{}
)

func InitRouter(r *gin.Engine) {
	initGetRouter(r)
	initPostRouter(r)
	initPatchRouter(r)
	initDeleteRouter(r)
}

func initGetRouter(r *gin.Engine) {
	for path, f := range getRouter {
		r.GET(path, f...)
	}
}

func initPostRouter(r *gin.Engine) {
	for path, f := range getRouter {
		r.POST(path, f...)
	}
}

func initPatchRouter(r *gin.Engine) {
	for path, f := range getRouter {
		r.PATCH(path, f...)
	}
}

func initDeleteRouter(r *gin.Engine) {
	for path, f := range getRouter {
		r.DELETE(path, f...)
	}
}
