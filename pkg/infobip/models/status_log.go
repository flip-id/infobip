package models

import (
	"github.com/jinzhu/gorm"
)

type StatusLog struct {
	gorm.Model
	BulkId      string `gorm:"type:varchar(50);null;index"`
	MessageId   string `gorm:"type:varchar(50);not null;index"`
	PhoneNumber string `gorm:"type:varchar(20);not null;index"`
	StatusCode  string `gorm:"type:varchar(50);not null;index"`
}
