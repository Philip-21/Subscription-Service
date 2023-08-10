package database

import (
	"gorm.io/gorm"
)

var db *gorm.DB

// New is the function used to create an instance of the data package. It returns the type
// Model, which embeds all the types we want to be available to our application.
func New(dbPool *gorm.DB) Models {
	db = dbPool

	return Models{
		User: &User{},
		Plan: &Plan{},
	}
}

// Models is the type for this package. Note that any model that is included as a member
// in this type is available to us throughout the application, anywhere that the
// app variable is used, provided that the model is also added in the New function.
type Models struct {
	User UserInterface
	Plan PlanInterface
}
