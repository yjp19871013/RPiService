package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/yjp19871013/RPiService/filestation/download_proxy"

	"github.com/gin-gonic/gin"
	"github.com/yjp19871013/RPiService/filestation/dto"
)

func AddDownloadFile(c *gin.Context) {
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

	id, err := download_proxy.GetInstance().AddTask(request.Url, request.SaveFilename)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := dto.DownloadFileResponse{
		ID:           id,
		Url:          request.Url,
		SaveFilename: request.SaveFilename,
	}
	c.JSON(http.StatusOK, response)
}

func DeleteDownloadFile(c *gin.Context) {
	var request dto.DeleteDownloadFileRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = download_proxy.GetInstance().DeleteTask(request.ID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := dto.DeleteDownloadFileResponse{
		ID: request.ID,
	}
	c.JSON(http.StatusOK, response)
}

func DownloadProgressPush(c *gin.Context) {
	var request dto.DownloadProgressRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
}
