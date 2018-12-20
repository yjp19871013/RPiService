package download_proxy

import (
	"fmt"
	"sync"

	"github.com/yjp19871013/RPiService/filestation/db"
)

const (
	ErrAlreadyExist = "already exist"
)

type DownloadProxy struct {
	sync.Mutex
}

func NewDownloadProxy() *DownloadProxy {
	return &DownloadProxy{}
}

func (proxy *DownloadProxy) AddDownloadTask(urlStr string, saveFilename string) (uint, error) {
	proxy.Lock()
	defer proxy.Unlock()

	var downloadTask db.DownloadTask
	err := db.GetInstance().Where("url = ?", urlStr).First(&downloadTask).Error
	if err == nil {
		return 0, fmt.Errorf(ErrAlreadyExist)
	}

	downloadTask.Url = urlStr
	downloadTask.SaveFileName = saveFilename

	err = db.GetInstance().Save(&downloadTask).Error
	if err != nil {
		return 0, err
	}

	return downloadTask.ID, nil
}

func (proxy *DownloadProxy) DeleteDownloadTask(urlStr string) (uint, error) {
	proxy.Lock()
	defer proxy.Unlock()

	var downloadTask db.DownloadTask
	err := db.GetInstance().Where("url = ?", urlStr).First(&downloadTask).Error
	if err != nil {
		return 0, err
	}

	err = db.GetInstance().Delete(&downloadTask).Error
	if err != nil {
		return 0, err
	}

	return downloadTask.ID, nil
}
