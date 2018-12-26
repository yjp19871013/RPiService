package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yjp19871013/RPiService/filestation/api"
)

var (
	getRouter = map[string][]gin.HandlerFunc{
		"/api/file-station/download-proxy/files/download-progresses/:ids": { /*middleware.JWTValidateMiddleware(), */ api.DownloadProgresses},
	}

	postRouter = map[string][]gin.HandlerFunc{
		"/api/file-station/download-proxy/files": { /*middleware.JWTValidateMiddleware(), */ api.AddDownloadFile},
	}

	patchRouter = map[string][]gin.HandlerFunc{}

	deleteRouter = map[string][]gin.HandlerFunc{
		"/api/file-station/download-proxy/files/:id": { /*middleware.JWTValidateMiddleware(), */ api.DeleteDownloadFile},
	}
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
	for path, f := range postRouter {
		r.POST(path, f...)
	}
}

func initPatchRouter(r *gin.Engine) {
	for path, f := range patchRouter {
		r.PATCH(path, f...)
	}
}

func initDeleteRouter(r *gin.Engine) {
	for path, f := range deleteRouter {
		r.DELETE(path, f...)
	}
}
