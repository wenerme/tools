package models

import "gorm.io/datatypes"

type Setting struct {
	Model
	Name    string `gorm:"unique"`
	Value   datatypes.JSON
	Version string
}
