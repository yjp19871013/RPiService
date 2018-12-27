package db

type FileInfo struct {
	ID           uint   `gorm:"primary_key"`
	FilePathname string `gorm:"unique;not null"`
}

func SaveFileInfo(info *FileInfo) error {
	return db.Save(info).Error
}

func FindFileInfoByFilePathname(pathname string) (*FileInfo, error) {
	fileInfo := &FileInfo{}
	err := db.Where("file_pathname = ?", pathname).First(fileInfo).Error
	if err != nil {
		return nil, err
	}

	return fileInfo, nil
}
