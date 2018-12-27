package filestation

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/yjp19871013/RPiService/utils"

	"github.com/yjp19871013/RPiService/middleware"

	"github.com/yjp19871013/RPiService/api/filestation/download_proxy"
	"github.com/yjp19871013/RPiService/api/filestation/dto"

	"github.com/gin-gonic/gin"
	"github.com/yjp19871013/RPiService/db"
)

const (
	saveDir = "files/"
)

func GetDownloadTasks(c *gin.Context) {
	userContext := c.Value(middleware.ContextUserKey)
	if userContext == nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, _ := userContext.(*db.User)

	tasks, err := download_proxy.GetInstance().GetUserTasks(user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := dto.GetAllTaskResponse{
		Tasks: make([]dto.DownloadTask, 0),
	}

	for _, task := range tasks {
		response.Tasks = append(response.Tasks, dto.DownloadTask{
			ID:           task.ID,
			Url:          task.Url,
			SaveFilename: filepath.Base(task.SaveFilePathname),
		})
	}

	c.JSON(http.StatusOK, response)
}

func AddDownloadTask(c *gin.Context) {
	userContext := c.Value(middleware.ContextUserKey)
	if userContext == nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, _ := userContext.(*db.User)

	absSaveDir, err := filepath.Abs(user.Email + "/" + saveDir)
	if err != nil {
		panic("download proxy save dir abs error")
	}

	exist, err := utils.PathExists(absSaveDir)
	if err != nil {
		panic("download proxy save dir PathExists error")
	}

	if !exist {
		_ = os.MkdirAll(absSaveDir, os.ModeDir|os.ModePerm)
	}

	var request dto.AddDownloadTaskRequest
	err = c.ShouldBindJSON(&request)
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

	id, err := download_proxy.GetInstance().AddTask(request.Url, absSaveDir+"/"+saveFilename, user)
	if err != nil {
		if err == download_proxy.SavePathnameExistErr {
			c.AbortWithStatus(http.StatusConflict)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := dto.DownloadTask{
		ID:           id,
		Url:          request.Url,
		SaveFilename: saveFilename,
	}

	c.JSON(http.StatusOK, response)
}

func DeleteDownloadTask(c *gin.Context) {
	userContext := c.Value(middleware.ContextUserKey)
	if userContext == nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, _ := userContext.(*db.User)

	log.Println(user)

	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = download_proxy.GetInstance().DeleteTask(uint(idInt))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := dto.DeleteDownloadFileResponse{
		ID: uint(idInt),
	}
	c.JSON(http.StatusOK, response)
}

func DownloadTaskProgresses(c *gin.Context) {
	userContext := c.Value(middleware.ContextUserKey)
	if userContext == nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, _ := userContext.(*db.User)

	log.Println(user)

	ids := make([]uint, 0)
	idsStr := strings.Split(c.Param("ids"), ";")
	for _, id := range idsStr {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		ids = append(ids, uint(idInt))
	}

	progresses, err := download_proxy.GetInstance().GetProcesses(ids)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := &dto.GetDownloadProgressResponse{
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

func GetFiles(c *gin.Context) {
	userContext := c.Value(middleware.ContextUserKey)
	if userContext == nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, _ := userContext.(*db.User)

	log.Println(user)

	infos, err := db.FindAllFileInfos()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := dto.GetAllFileInfosResponse{
		FileInfos: make([]dto.FileInfo, 0),
	}
	for _, info := range infos {
		response.FileInfos = append(response.FileInfos, dto.FileInfo{
			ID:       info.ID,
			FileName: filepath.Base(info.FilePathname),
		})
	}

	c.JSON(http.StatusOK, response)
}
