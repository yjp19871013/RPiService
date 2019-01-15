package db

type Role struct {
	ID          uint   `gorm:"primary_key"`
	Name        string `gorm:"unique;not null"`
	Description string
	Permissions []Permission `gorm:"many2many:role_permissions"`
}

func FindRoleByName(name string) (*Role, error) {
	role := &Role{}
	err := db.Where("name = ?", name).First(role).Error
	if err != nil {
		return nil, err
	}

	return role, nil
}

func GetAllRoles() ([]Role, error) {
	roles := make([]Role, 0)
	err := db.Find(&roles).Error
	if err != nil {
		return nil, err
	}

	return roles, err
}

func SaveRole(role *Role) error {
	return db.Save(role).Error
}
