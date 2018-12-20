package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/yjp19871013/RPiService/filestation/download_proxy"

	"github.com/gin-gonic/gin"
	"github.com/yjp19871013/RPiService/filestation/dto"
)

var downloadProxy = download_proxy.NewDownloadProxy()

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

	task := download_proxy.Task{
		Url:          request.Url,
		SaveFilename: saveFilename,
	}

	id, err := downloadProxy.AddDownloadTask(task)
	if err != nil {
		if err.Error() == download_proxy.ErrAlreadyExist {
			c.AbortWithStatus(http.StatusConflict)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := dto.DownloadFileResponse{
		ID:           id,
		Url:          request.Url,
		SaveFilename: saveFilename,
	}
	c.JSON(http.StatusOK, response)
}
