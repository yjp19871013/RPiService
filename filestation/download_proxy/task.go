package download_proxy

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Task struct {
	sync.Mutex
	Url              string
	SaveFilePathname string
	progress         uint
	isStart          bool

	logFile      *os.File
	cmd          *exec.Cmd
	completeChan chan bool
}

func NewTask(urlStr string, saveFilePathname string, completeChan chan bool) *Task {
	f, err := ioutil.TempFile("/tmp", "download_proxy")
	if err != nil {
		return nil
	}

	return &Task{
		Url:              urlStr,
		SaveFilePathname: saveFilePathname,
		logFile:          f,
		completeChan:     completeChan,
	}
}

func (task *Task) Start() error {
	task.Lock()
	defer task.Unlock()

	if task.isStart {
		return nil
	}

	log.Println(task.logFile.Name())

	task.cmd = exec.Command("wget", "-c", task.Url, "-O", task.SaveFilePathname, "-o", task.logFile.Name(), "-v")
	err := task.cmd.Start()
	if err != nil {
		return err
	}

	task.isStart = true

	go task.ParseProgress()

	return nil
}

func (task *Task) Stop() error {
	task.Lock()
	defer task.Unlock()

	if !task.isStart {
		return nil
	}

	_ = task.cmd.Process.Kill()
	_ = task.cmd.Wait()

	task.isStart = false

	_ = task.logFile.Close()
	_ = os.Remove(task.logFile.Name())

	return nil
}

func (task *Task) GetProgress() uint {
	return task.progress
}

func (task *Task) ParseProgress() {
	for true {
		content, _ := ioutil.ReadAll(task.logFile)
		outputStr := string(content)

		endIndex := strings.LastIndex(outputStr, "%")
		if endIndex == -1 {
			continue
		}

		progress := strings.TrimSpace(outputStr[endIndex-3 : endIndex])
		progressInt, err := strconv.Atoi(progress)
		if err != nil {
			continue
		}

		task.progress = uint(progressInt)

		if task.progress == 100 {
			task.completeChan <- true
			return
		}

		time.Sleep(1 * time.Second)
	}
}
