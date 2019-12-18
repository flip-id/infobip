package database

import "github.com/flip-id/infobip/pkg/infobip/models"

func init() {
	InitMigration()
}

func InitMigration() {
	db, _ := Connect("root", "root", "infobip", "localhost")

	db.AutoMigrate(&models.StatusLog{})
}
