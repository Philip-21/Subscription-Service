package database

import (
	"time"
)

// UserInterface is the interface for the user type. In order
// to satisfy this interface, all specified methods must be implemented.
// We do this so we can test things easily. Both data.User and data.UserTest
// implement this interface.
type UserInterface interface {
	GetAll() ([]*User, error)
	GetByEmail(email string) (*User, error)
	GetOne(id int) (*User, error)
	Update(user User) error
	// Delete() error
	DeleteByID(id int) error
	Insert(user User) (int, error)
	ResetPassword(password string) error
	PasswordMatches(plainText string) (bool, error)
}

// PlanInterface is the type for the plan type. Both data.Plan and data.PlanTest
// implement this interface.
type PlanInterface interface {
	GetAll() ([]*Plan, error)
	GetOne(id int) (*Plan, error)
	SubscribeUserToPlan(user User, plan Plan) error
	AmountForDisplay() string
}

// TestNew is the function used to create an instance of the data package. It returns the type
// Model, which embeds all the types we want to be available to our application. This
// is only used when running tests.
// func TestNew(dbPool *sql.DB) Models {
// 	db = dbPool

// 	return Models{
// 		User: UserTest{},
// 		Plan: PlanTest{},
// 	}
// }

// UserTest is the structure which holds one user from the database,
// and is used for testing.
type UserTest struct {
	ID        int
	Email     string
	FirstName string
	LastName  string
	Password  string
	Active    int
	IsAdmin   int
	CreatedAt time.Time
	UpdatedAt time.Time
	Plan      *Plan
}

type PlanTest struct {
	ID                  int
	PlanName            string
	PlanAmount          int
	PlanAmountFormatted string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
