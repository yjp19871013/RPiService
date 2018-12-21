package download_proxy

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type Task struct {
	Url              string `gorm:"unique;not null"`
	SaveFilePathname string `gorm:"not null"`
	Process          uint
	cmd              *exec.Cmd
}

func (task *Task) Start() error {
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
			fmt.Println(outputStr)
			endIndex := strings.Index(outputStr, "%")
			if endIndex == -1 {
				continue
			}

			process := strings.TrimSpace(outputStr[endIndex-3 : endIndex])
			processInt, err := strconv.Atoi(process)
			if err != nil {
				continue
			}

			task.Process = uint(processInt)
		}
	}(reader)

	return nil
}

func (task *Task) Stop() error {
	return task.cmd.Process.Kill()
}
