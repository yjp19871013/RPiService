package download_proxy

import (
	"bufio"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

type Task struct {
	sync.Mutex
	Url              string
	SaveFilePathname string
	Progress         uint
	IsStart          bool
	cmd              *exec.Cmd
}

func (task *Task) Start() error {
	task.Lock()
	defer task.Unlock()

	if task.IsStart {
		return nil
	}

	task.cmd = exec.Command("wget", "-c", task.Url, "-O", task.SaveFilePathname, "-o", "/dev/stdout")
	stdout, err := task.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	err = task.cmd.Start()
	if err != nil {
		return err
	}

	reader := bufio.NewReader(stdout)
	go func(reader *bufio.Reader) {
		for true {
			line, _, _ := reader.ReadLine()
			outputStr := string(line)
			endIndex := strings.Index(outputStr, "%")
			if endIndex == -1 {
				continue
			}

			progress := strings.TrimSpace(outputStr[endIndex-3 : endIndex])
			progressInt, err := strconv.Atoi(progress)
			if err != nil {
				continue
			}

			task.Progress = uint(progressInt)
		}
	}(reader)

	task.IsStart = true

	return nil
}

func (task *Task) Stop() error {
	task.Lock()
	defer task.Unlock()

	if !task.IsStart {
		return nil
	}

	err := task.cmd.Process.Kill()
	if err != nil {
		return err
	}

	_ = task.cmd.Wait()

	task.IsStart = false

	return nil
}
