package datasources

import "gorm.io/gorm"

type Datasource struct {
	gorm.Model
	Name string `json:"name,omitempty" gorm:"not null"`
}
