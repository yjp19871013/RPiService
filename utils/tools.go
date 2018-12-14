package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func LoadJsonFileConfig(pathname string, v interface{}) error {
	absPath, err := filepath.Abs(pathname)
	if err != nil {
		fmt.Println(err)
		return err
	}

	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func GetAllFiles(pathname string) ([]string, error) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return nil, err
	}

	files := make([]string, 0)
	for _, fi := range rd {
		if fi.IsDir() {
			GetAllFiles(pathname + fi.Name() + "\\")
		} else {
			files = append(files, fi.Name())
		}
	}

	return files, err
}
