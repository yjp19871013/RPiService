package db

type DownloadTask struct {
	ID           uint   `gorm:"primary_key"`
	Url          string `gorm:"not null"`
	SaveFileName string `gorm:"not null"`
}

func SaveDownloadTask(task *DownloadTask) (*DownloadTask, error) {
	err := db.Save(task).Error
	if err != nil {
		return nil, err
	}

	return task, nil
}
