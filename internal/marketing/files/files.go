package files

import "gorm.io/gorm"

type File struct {
	gorm.Model
	Name    string `json:"name,omitempty" gorm:"not null"`
	Success bool   `json:"success,omitempty" gorm:"not null"`
}
