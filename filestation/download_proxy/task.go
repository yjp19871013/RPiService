package download_proxy

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

type Task struct {
	sync.Mutex
	Url              string
	SaveFilePathname string
	progress         uint
	isStart          bool

	logFile *os.File
	cmd     *exec.Cmd
}

func NewTask(urlStr string, saveFilePathname string) *Task {
	f, err := ioutil.TempFile("/tmp", "download_proxy")
	if err != nil {
		return nil
	}

	return &Task{
		Url:              urlStr,
		SaveFilePathname: saveFilePathname,
		logFile:          f,
	}
}

func (task *Task) Start() error {
	task.Lock()
	defer task.Unlock()

	if task.isStart {
		return nil
	}

	fmt.Println(task.logFile.Name())

	task.cmd = exec.Command("wget", "-c", task.Url, "-O", task.SaveFilePathname, "-o", task.logFile.Name())
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
		if strings.Count(outputStr, "%") == 0 {
			task.progress = 100
			break
		}

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
	}
}
