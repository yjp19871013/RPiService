package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/mail"
	"net/smtp"
	"os"
	"path/filepath"
	"time"

	"github.com/yjp19871013/RPiService/settings"

	"github.com/scorredoira/email"

	"github.com/json-iterator/go"
)

func JsonMarshal(v interface{}) (string, error) {
	j := jsoniter.ConfigCompatibleWithStandardLibrary
	return j.MarshalToString(v)
}

func JsonUnmarshal(str string, v interface{}) error {
	j := jsoniter.ConfigCompatibleWithStandardLibrary
	return j.UnmarshalFromString(str, v)
}

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

	err = JsonUnmarshal(string(data), v)
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

func FileSize(path string) (uint64, error) {
	stat, err := os.Stat(path)
	if err == nil {
		return uint64(stat.Size()), nil
	}

	if os.IsNotExist(err) {
		return 0, nil
	}

	return 0, err
}

func GetAllFiles(pathname string) ([]string, error) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return nil, err
	}

	files := make([]string, 0)
	for _, fi := range rd {
		if fi.IsDir() {
			GetAllFiles(pathname + fi.Name() + string(os.PathSeparator))
		} else {
			files = append(files, fi.Name())
		}
	}

	return files, err
}

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func GenerateValidateCode() string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	return fmt.Sprintf("%06d", r.Intn(1000000))
}

func SendEmail(subject string, body string, to string) error {
	m := email.NewMessage(subject, body)
	m.From = mail.Address{Name: "From", Address: settings.SMTPUsername}
	m.To = []string{to}

	auth := smtp.PlainAuth("", settings.SMTPUsername, settings.SMTPPassword, settings.SMTPHost)
	return email.Send(settings.SMTPAddress, auth, m)
}
