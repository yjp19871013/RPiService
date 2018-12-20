package db

type DownloadTask struct {
	ID           uint   `gorm:"primary_key"`
	Url          string `gorm:"not null"`
	SaveFileName string `gorm:"not null"`
}
