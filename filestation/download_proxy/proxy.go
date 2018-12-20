package download_proxy

type Proxy struct {
}

func NewProxy() *Proxy {
	return &Proxy{}
}

func (proxy *Proxy) AddTask(urlStr string, saveFilename string) error {
	return nil
}
