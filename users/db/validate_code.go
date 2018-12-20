package db

type ValidateCode struct {
	ID           uint   `gorm:"primary_key"`
	Email        string `gorm:"unique;not null"`
	ValidateCode string `gorm:"not null"`
}

func FindValidateCodeByEmail(email string) (*ValidateCode, error) {
	validateCode := &ValidateCode{}
	err := db.Where("email = ?", email).First(validateCode).Error
	if err != nil {
		return nil, err
	}

	return validateCode, nil
}

func SaveValidateCode(validateCode *ValidateCode) error {
	return db.Save(validateCode).Error
}

func DeleteValidateCode(validateCode *ValidateCode) error {
	return db.Delete(validateCode).Error
}
