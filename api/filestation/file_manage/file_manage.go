package file_manage

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/yjp19871013/RPiService/settings"

	"github.com/gin-gonic/gin"
	"github.com/yjp19871013/RPiService/api/filestation/file_manage/dto"
	"github.com/yjp19871013/RPiService/db"
	"github.com/yjp19871013/RPiService/middleware"
)

const (
	FilesRelativeDir = "files/"
)

func GetFiles(c *gin.Context) {
	userContext := c.Value(middleware.ContextUserKey)
	if userContext == nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, _ := userContext.(*db.User)

	infos, err := db.FindFileInfosByUser(user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := dto.GetAllFileInfosResponse{
		FileInfos: make([]dto.FileInfo, 0),
	}
	for _, info := range infos {
		response.FileInfos = append(response.FileInfos, dto.FileInfo{
			ID:           info.ID,
			FileName:     filepath.Base(info.FilePathname),
			CompleteDate: info.CompleteDate,
			SizeKb:       info.SizeKb,
		})
	}

	c.JSON(http.StatusOK, response)
}

func DownloadFile(c *gin.Context) {
	userContext := c.Value(middleware.ContextUserKey)
	if userContext == nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, _ := userContext.(*db.User)

	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	fileInfo, err := db.FindFileInfoById(uint(idInt))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if fileInfo.UserId != user.ID {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	staticUrl := settings.ServerUrl + settings.StaticRoot + "/" +
		FilesRelativeDir + user.Email + "/" + filepath.Base(fileInfo.FilePathname)
	response := dto.DownloadFileResponse{
		StaticUrl: staticUrl,
	}

	c.JSON(http.StatusOK, response)
}
