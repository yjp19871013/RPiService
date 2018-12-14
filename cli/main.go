package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/yjp19871013/RPiService/utils"

	DEATH "gopkg.in/vrecan/death.v3"
)

var (
	start = flag.Bool("start", false, "Start daemon, start")
	stop  = flag.Bool("stop", false, "Stop daemon")
)

// 检查错误，如果有错误就退出
func checkErrorAndExit(err error, exitCode int) {
	if err != nil {
		fmt.Println(err)
		os.Exit(exitCode)
	}
}

// 参数是否有效
func isValidArgs(args []string) bool {
	return len(args) == 0 || (len(args) == 1 && *start) || (len(args) == 1 && *stop)
}

func main() {
	flag.Parse()

	if !isValidArgs(flag.Args()) {
		fmt.Println("Usage error, use -help to see how to use")
		os.Exit(1)
	}

	// 同时指定了-start和-stop
	if *start && *stop {
		fmt.Println("You can only use start or stop.")
		os.Exit(1)
	}

	config := &Config{}
	err := utils.LoadJsonFileConfig("cli_config.json", config)
	checkErrorAndExit(err, 1)

	exist, err := utils.PathExists(config.PIDFileDir)
	checkErrorAndExit(err, 1)

	if !exist {
		pidFileDirPath, err := filepath.Abs(config.PIDFileDir)
		checkErrorAndExit(err, 1)

		err = os.MkdirAll(pidFileDirPath, 0777)
		checkErrorAndExit(err, 1)
	}

	commands := make([]*exec.Cmd, 0)
	for _, pathname := range config.Processes {
		commands = append(commands, exec.Command(pathname))
	}

	if !*stop && !*start {
		// 没有指定-start或者-stop，阻塞启动程序
		for _, cmd := range commands {
			cmd.Start()
		}

		death := DEATH.NewDeath(syscall.SIGINT, syscall.SIGTERM)
		death.WaitForDeath()

		os.Exit(0)
	} else if *start && !*stop {
		// 指定了-start，启动为守护进程
		if os.Getppid() != 1 {
			for _, cmd := range commands {
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				err := cmd.Start()
				if err != nil {
					panic(err)
				}

				ioutil.WriteFile(config.PIDFileDir+string(filepath.Separator)+strconv.Itoa(cmd.Process.Pid)+".pid",
					[]byte(strconv.Itoa(cmd.Process.Pid)), 0777)
			}
		}

		os.Exit(0)
	} else if !*start && *stop {
		// 指定了-stop，关闭相关的进程
		pidFiles, err := utils.GetAllFiles(config.PIDFileDir)
		checkErrorAndExit(err, 1)

		for _, file := range pidFiles {
			filePathname := config.PIDFileDir + string(filepath.Separator) + file
			pid, err := ioutil.ReadFile(filePathname)
			checkErrorAndExit(err, 1)

			pidInt, err := strconv.Atoi(string(pid))
			checkErrorAndExit(err, 1)

			process, err := os.FindProcess(pidInt)
			if err != nil {
				os.Remove(filePathname)
				continue
			}

			process.Kill()

			os.Remove(filePathname)
		}
	}
}
