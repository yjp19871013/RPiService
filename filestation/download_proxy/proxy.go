package download_proxy

import (
	"errors"
	"log"
	"os"
	"sync"

	"github.com/yjp19871013/RPiService/filestation/db"
)

var (
	SavePathnameExistErr = errors.New("Save pathname has exist")

	downloadProxy *Proxy
)

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

func (proxy *Proxy) GetAllTasks() ([]db.DownloadTask, error) {
	proxy.Lock()
	defer proxy.Unlock()

	tasks, err := db.FindAllDownloadTasks()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (proxy *Proxy) AddTask(urlStr string, saveFilePathname string) (uint, error) {
	proxy.Lock()
	defer proxy.Unlock()

	downloadTask, err := db.FindDownloadTaskByUrl(urlStr)
	if err == nil {
		return downloadTask.ID, nil
	}

	downloadTask, err = db.FindDownloadTaskBySaveFilePathname(saveFilePathname)
	if err == nil {
		return 0, SavePathnameExistErr
	}

	downloadTask = &db.DownloadTask{
		Url:              urlStr,
		SaveFilePathname: saveFilePathname,
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

func (proxy *Proxy) GetProcesses(ids []uint) (map[uint]uint, error) {
	proxy.Lock()
	defer proxy.Unlock()

	progresses := make(map[uint]uint, 0)
	for _, id := range ids {
		task := proxy.taskMap[id]
		if task != nil {
			progress := task.GetProgress()
			progresses[id] = progress
		} else {
			progresses[id] = 100
		}
	}

	return progresses, nil
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
	err := db.SaveDownloadTask(downloadTask)
	if err != nil {
		return err
	}

	completeChan := make(chan bool)
	task := NewTask(downloadTask.Url, downloadTask.SaveFilePathname, completeChan)

	go func(id uint, completeChan chan bool) {
		for true {
			select {
			case complete := <-completeChan:
				proxy.Lock()

				if complete {
					err := db.DeleteDownloadTask(&db.DownloadTask{ID: id})
					if err != nil {
						proxy.Unlock()
						continue
					}

					_ = task.Stop()

					delete(proxy.taskMap, id)

					proxy.Unlock()
					return
				}

				proxy.Unlock()
			}
		}
	}(downloadTask.ID, completeChan)

	err = task.Start()
	if err != nil {
		return err
	}

	proxy.taskMap[downloadTask.ID] = task

	return nil
}
