package download_proxy

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type Proxy struct {
	sync.Mutex
	taskMap map[string]*Task
}

func NewProxy() *Proxy {
	return &Proxy{
		taskMap: make(map[string]*Task),
	}
}

func (proxy *Proxy) Start(initTasks map[string]string) error {
	proxy.Lock()
	defer proxy.Unlock()

	for urlStr, saveFilePathname := range initTasks {
		err := proxy.addTaskWithoutLock(urlStr, saveFilePathname)
		if err != nil {
			log.Println("Start task", urlStr, "error")
			return err
		}
	}

	return nil
}

func (proxy *Proxy) Stop() error {
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

func (proxy *Proxy) AddTask(urlStr string, saveFilePathname string) error {
	proxy.Lock()
	defer proxy.Unlock()

	err := proxy.addTaskWithoutLock(urlStr, saveFilePathname)
	if err != nil {
		return err
	}

	return nil
}

func (proxy *Proxy) DeleteTask(urlStr string) error {
	task := proxy.taskMap[urlStr]
	if task == nil {
		return nil
	}

	err := task.Stop()
	if err != nil {
		return err
	}

	err = os.Remove(task.SaveFilePathname)
	if err != nil {
		return err
	}

	delete(proxy.taskMap, task.Url)

	return nil
}

func (proxy *Proxy) GetProcess(urlStr string) (uint, error) {
	proxy.Lock()
	defer proxy.Unlock()

	task := proxy.taskMap[urlStr]
	if task == nil {
		return 0, fmt.Errorf("No this task")
	}

	return task.Process, nil
}

func (proxy *Proxy) addTaskWithoutLock(urlStr string, saveFilePathname string) error {
	task := Task{
		Url:              urlStr,
		SaveFilePathname: saveFilePathname,
	}

	err := task.Start()
	if err != nil {
		return err
	}

	proxy.taskMap[task.Url] = &task

	return nil
}
