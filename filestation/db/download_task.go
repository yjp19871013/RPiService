package db

import "log"

type DownloadTask struct {
	ID               uint   `gorm:"primary_key"`
	Url              string `gorm:"unique;not null"`
	SaveFilePathname string `gorm:"unique;not null"`
}

func FindAllDownloadTasks() ([]DownloadTask, error) {
	tasks := make([]DownloadTask, 0)
	err := db.Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func FindDownloadTaskByUrl(url string) (*DownloadTask, error) {
	task := &DownloadTask{}
	err := db.Where("url = ?", url).First(task).Error
	if err != nil {
		return nil, err
	}

	return task, nil
}

func FindDownloadTaskBySaveFilePathname(savePathname string) (*DownloadTask, error) {
	task := &DownloadTask{}
	err := db.Where("save_file_pathname = ?", savePathname).First(task).Error
	if err != nil {
		return nil, err
	}

	log.Println(task)
	return task, nil
}

func SaveDownloadTask(task *DownloadTask) error {
	return db.Save(task).Error
}

func DeleteDownloadTask(task *DownloadTask) error {
	return db.Delete(task).Error
}
