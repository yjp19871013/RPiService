package filestation

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/yjp19871013/RPiService/api/filestation/dto"
	"github.com/yjp19871013/RPiService/db"
	"github.com/yjp19871013/RPiService/middleware"
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

	log.Println(infos)

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
