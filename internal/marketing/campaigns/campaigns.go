package campaigns

import "gorm.io/gorm"

type Campaign struct {
	gorm.Model
	Name string `json:"name,omitempty" gorm:"not null"`
}
