package download_proxy

const (
	wgetCmd = "wget %s -o /dev/stdout -O %s"
)

type Task struct {
	Url          string
	SaveFilename string
}

func (task *Task) Start() {
	go func() {
		//cmd := exec.Command(fmt.Sprintf(wgetCmd, task.Url, task.SaveFilename))
	}()
}
