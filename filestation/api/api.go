package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/yjp19871013/RPiService/filestation/db"

	"github.com/yjp19871013/RPiService/filestation/download_proxy"

	"github.com/gin-gonic/gin"
	"github.com/yjp19871013/RPiService/filestation/dto"
)

var downloadProxy = download_proxy.NewProxy()

func DownloadFile(c *gin.Context) {
	var request dto.DownloadFileRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	saveFilename := request.SaveFilename
	if len(saveFilename) == 0 {
		saveFilename = request.Url[strings.LastIndex(request.Url, "/")+1:]
	}

	err = downloadProxy.AddTask(request.Url, saveFilename)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	downloadTask := &db.DownloadTask{}
	downloadTask, err = db.SaveDownloadTask(downloadTask)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := dto.DownloadFileResponse{
		ID:           downloadTask.ID,
		Url:          request.Url,
		SaveFilename: saveFilename,
	}

	c.JSON(http.StatusOK, response)
}
