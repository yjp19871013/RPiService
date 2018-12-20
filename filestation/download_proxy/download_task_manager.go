package download_proxy

import (
	"fmt"
	"sync"

	"github.com/yjp19871013/RPiService/filestation/db"
)

const (
	ErrAlreadyExist = "already exist"
)

type Task struct {
	Url          string
	SaveFilename string
}

type DownloadProxy struct {
	sync.Mutex
}

func NewDownloadProxy() *DownloadProxy {
	return &DownloadProxy{}
}

func (proxy *DownloadProxy) AddDownloadTask(task Task) (uint, error) {
	proxy.Lock()
	defer proxy.Unlock()

	var downloadTask db.DownloadTask
	err := db.GetInstance().Where("url = ?", task.Url).First(&downloadTask).Error
	if err == nil {
		return 0, fmt.Errorf(ErrAlreadyExist)
	}

	downloadTask.Url = task.Url
	downloadTask.SaveFileName = task.SaveFilename

	err = db.GetInstance().Save(&downloadTask).Error
	if err != nil {
		return 0, err
	}

	return downloadTask.ID, nil
}

func (proxy *DownloadProxy) DeleteDownloadTask(task Task) (uint, error) {
	proxy.Lock()
	defer proxy.Unlock()

	var downloadTask db.DownloadTask
	err := db.GetInstance().Where("url = ?", task.Url).First(&downloadTask).Error
	if err != nil {
		return 0, err
	}

	err = db.GetInstance().Delete(&downloadTask).Error
	if err != nil {
		return 0, err
	}

	return downloadTask.ID, nil
}
