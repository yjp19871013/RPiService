package api

import (
	"log"
	"net/http"
	"strconv"
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
		startIndex := strings.LastIndex(request.Url, "/") + 1
		endIndex := strings.LastIndex(request.Url, "?")
		if endIndex != -1 {
			saveFilename = request.Url[startIndex:endIndex]
		} else {
			saveFilename = request.Url[startIndex:]
		}
	}

	id, err := download_proxy.GetInstance().AddTask(request.Url, saveFilename)
	if err != nil {
		if err == download_proxy.SavePathnameExistErr {
			c.AbortWithStatus(http.StatusConflict)
			return
		}

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

func DownloadProgresses(c *gin.Context) {
	ids := make([]uint, 0)
	idsStr := strings.Split(c.Param("ids"), ";")
	for _, id := range idsStr {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		}

		ids = append(ids, uint(idInt))
	}

	progresses, err := download_proxy.GetInstance().GetProcesses(ids)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	response := &dto.DownloadProgressResponse{
		Progresses: make([]dto.DownloadProgress, 0),
	}

	for _, id := range ids {
		response.Progresses = append(response.Progresses, dto.DownloadProgress{
			ID:       id,
			Progress: progresses[id],
		})
	}

	c.JSON(http.StatusOK, response)
}
