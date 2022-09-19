package models

import "gorm.io/gorm"

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&Artist{},
		&Track{},
	)

	if err != nil {
		panic(err)
	}
}
