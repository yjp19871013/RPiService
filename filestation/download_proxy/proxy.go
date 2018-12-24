package download_proxy

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/yjp19871013/RPiService/filestation/db"
)

var (
	AlreadyExistErr = errors.New("Already exits")
	StartTaskErr    = errors.New("Start task error")
)

var downloadProxy *Proxy

func StartProxy() {
	downloadProxy = NewProxy()
	err := downloadProxy.start()
	if err != nil {
		panic("Download Proxy start error")
	}
}

func StopProxy() {
	err := downloadProxy.stop()
	if err != nil {
		panic("Download Proxy stop error")
	}
}

func GetInstance() *Proxy {
	return downloadProxy
}

type Proxy struct {
	sync.Mutex
	taskMap map[uint]*Task
}

func NewProxy() *Proxy {
	return &Proxy{
		taskMap: make(map[uint]*Task),
	}
}

func (proxy *Proxy) AddTask(urlStr string, saveFilePathname string) (uint, error) {
	proxy.Lock()
	defer proxy.Unlock()

	downloadTask, err := db.FindDownloadTaskByUrl(urlStr)
	if err != nil {
		downloadTask = &db.DownloadTask{
			Url:              urlStr,
			SaveFilePathname: saveFilePathname,
		}
	}

	err = proxy.addTaskWithoutLock(downloadTask)
	if err != nil {
		return 0, err
	}

	return downloadTask.ID, nil
}

func (proxy *Proxy) DeleteTask(id uint) error {
	proxy.Lock()
	defer proxy.Unlock()

	err := db.DeleteDownloadTask(&db.DownloadTask{ID: id})
	if err != nil {
		return err
	}

	task := proxy.taskMap[id]
	if task == nil {
		return nil
	}

	err = task.Stop()
	if err != nil {
		return err
	}

	err = os.Remove(task.SaveFilePathname)
	if err != nil {
		return err
	}

	delete(proxy.taskMap, id)

	return nil
}

func (proxy *Proxy) GetProgress(id uint) (uint, error) {
	proxy.Lock()
	defer proxy.Unlock()

	task := proxy.taskMap[id]
	if task == nil {
		return 0, fmt.Errorf("No this task")
	}

	return task.Progress, nil
}

func (proxy *Proxy) start() error {
	proxy.Lock()
	defer proxy.Unlock()

	initTasks, err := db.FindAllDownloadTasks()
	if err != nil {
		return err
	}

	for _, downloadTask := range initTasks {
		err := proxy.addTaskWithoutLock(&downloadTask)
		if err != nil {
			log.Println("Start task", downloadTask.ID, downloadTask.Url, "error")
			return err
		}
	}

	return nil
}

func (proxy *Proxy) stop() error {
	proxy.Lock()
	defer proxy.Unlock()

	for id, task := range proxy.taskMap {
		err := task.Stop()
		if err != nil {
			log.Println("Stop task", task.Url, "error")
			continue
		}

		delete(proxy.taskMap, id)
	}

	return nil
}

func (proxy *Proxy) addTaskWithoutLock(downloadTask *db.DownloadTask) error {
	if downloadTask.ID == 0 {
		err := db.SaveDownloadTask(downloadTask)
		if err != nil {
			return err
		}
	}

	task := Task{
		Url:              downloadTask.Url,
		SaveFilePathname: downloadTask.SaveFilePathname,
	}

	err := task.Start()
	if err != nil {
		return StartTaskErr
	}

	proxy.taskMap[downloadTask.ID] = &task

	return nil
}
