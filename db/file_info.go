package db

type FileInfo struct {
	ID           uint    `gorm:"primary_key"`
	FilePathname string  `gorm:"unique;not null"`
	CompleteDate string  `gorm:"not null"`
	SizeKb       float64 `gorm:"not null"`

	UserId uint
}

func SaveFileInfo(info *FileInfo) error {
	return db.Save(info).Error
}

func FindFileInfosByUser(user *User) ([]FileInfo, error) {
	infos := make([]FileInfo, 0)
	err := db.Model(user).Related(&infos).Error
	if err != nil {
		return nil, err
	}

	return infos, nil
}

func FindFileInfoByFilePathname(pathname string) (*FileInfo, error) {
	fileInfo := &FileInfo{}
	err := db.Where("file_pathname = ?", pathname).First(fileInfo).Error
	if err != nil {
		return nil, err
	}

	return fileInfo, nil
}

func FindFileInfoById(id uint) (*FileInfo, error) {
	fileInfo := &FileInfo{}
	err := db.Where("id = ?", id).First(fileInfo).Error
	if err != nil {
		return nil, err
	}

	return fileInfo, nil
}

func DeleteFileInfo(info *FileInfo) error {
	return db.Delete(info).Error
}
