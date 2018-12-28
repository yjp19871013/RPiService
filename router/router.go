package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yjp19871013/RPiService/api/filestation/download_proxy"
	"github.com/yjp19871013/RPiService/api/filestation/file_manage"
	"github.com/yjp19871013/RPiService/api/users"
	"github.com/yjp19871013/RPiService/middleware"
	"github.com/yjp19871013/RPiService/settings"
)

var (
	getRouter = map[string][]gin.HandlerFunc{
		"/api/file-station/download-proxy/tasks":                          {middleware.JWTValidateMiddleware(), download_proxy.GetDownloadTasks},
		"/api/file-station/download-proxy/tasks/download-progresses/:ids": {middleware.JWTValidateMiddleware(), download_proxy.DownloadTaskProgresses},
		"/api/file-station/download-proxy/file-infos":                     {middleware.JWTValidateMiddleware(), file_manage.GetFiles},
		"/api/file-station/download-proxy/files/:id":                      {middleware.JWTValidateMiddleware(), file_manage.DownloadFile},
	}

	postRouter = map[string][]gin.HandlerFunc{
		"/api/users/token":         {users.CreateToken},
		"/api/users":               {users.Register},
		"/api/users/validate-code": {users.GenerateValidateCode},

		"/api/file-station/download-proxy/tasks": {middleware.JWTValidateMiddleware(), download_proxy.AddDownloadTask},
	}

	patchRouter = map[string][]gin.HandlerFunc{}

	deleteRouter = map[string][]gin.HandlerFunc{
		"/api/users/token": {middleware.JWTValidateMiddleware(), users.DeleteToken},

		"/api/file-station/download-proxy/tasks/:id": {middleware.JWTValidateMiddleware(), download_proxy.DeleteDownloadTask},
	}
)

func InitRouter(r *gin.Engine) {
	initGetRouter(r)
	initPostRouter(r)
	initPatchRouter(r)
	initDeleteRouter(r)

	r.Static(settings.StaticRoot, settings.StaticDir)
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
