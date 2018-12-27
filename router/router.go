package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yjp19871013/RPiService/api/filestation"
	"github.com/yjp19871013/RPiService/api/users"
	"github.com/yjp19871013/RPiService/middleware"
)

var (
	getRouter = map[string][]gin.HandlerFunc{
		"/api/file-station/download-proxy/tasks":                          {middleware.JWTValidateMiddleware(), filestation.GetDownloadTasks},
		"/api/file-station/download-proxy/tasks/download-progresses/:ids": {middleware.JWTValidateMiddleware(), filestation.DownloadTaskProgresses},
		"/api/file-station/download-proxy/files":                          {middleware.JWTValidateMiddleware(), filestation.GetFiles},
	}

	postRouter = map[string][]gin.HandlerFunc{
		"/api/users/token":         {users.CreateToken},
		"/api/users":               {users.Register},
		"/api/users/validate-code": {users.GenerateValidateCode},

		"/api/file-station/download-proxy/tasks": {middleware.JWTValidateMiddleware(), filestation.AddDownloadTask},
	}

	patchRouter = map[string][]gin.HandlerFunc{}

	deleteRouter = map[string][]gin.HandlerFunc{
		"/api/users/token": {middleware.JWTValidateMiddleware(), users.DeleteToken},

		"/api/file-station/download-proxy/tasks/:id": {middleware.JWTValidateMiddleware(), filestation.DeleteDownloadTask},
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
