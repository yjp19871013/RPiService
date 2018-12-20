package db

type Permission struct {
	ID          uint   `gorm:"primary_key"`
	Name        string `gorm:"unique;not null"`
	Description string
}

func FindPermissionByName(name string) (*Permission, error) {
	permission := &Permission{}
	err := db.Where("name = ?", name).First(permission).Error
	if err != nil {
		return nil, err
	}

	return permission, nil
}

func SavePermission(permission *Permission) error {
	return db.Save(permission).Error
}
