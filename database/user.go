package database

import (
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User is the structure which holds one user from the database.
type User struct {
	gorm.Model //had ID, created,updated at inplace
	Email     string
	FirstName string
	LastName  string
	Password  string
	Active    int
	IsAdmin   int
	Plan      *Plan     `gorm:"foreignKey:UserID"`
	UserPlan  *UserPlan `gorm:"foreignKey:UserID"`
}
type UserPlan struct {
	gorm.Model
	UserID uint
	PlanID uint
}

// GetAll returns a slice of all users, sorted by last name
func (u *User) GetAll() ([]*User, error) {
	var users []*User
	if err := db.Order("last_name").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetByEmail returns one user by email
func (u *User) GetByEmail(email string) (*User, error) {
	var user User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	// get plan, if any
	// query = `select p.id, p.plan_name, p.plan_amount, p.created_at, p.updated_at from
	// 		plans p
	// 		left join user_plans up on (p.id = up.plan_id)
	// 		where up.user_id = $1`
	var plan Plan
	result = db.Joins("left join userplan up on plans.id = up.plan_id").
		Where("up.user_id = ?", user.ID).
		First(&plan)
	if result.Error == nil {
		user.Plan = &plan
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return &user, nil
}

// GetOne returns one user by id
func (u *User) GetOne(id uint) (*User, error) {
	var user User
	if err := db.Where("id = ?", id).Preload("plan").First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates one user in the database, using the information
// stored in the receiver u
func (u *User) Update(user User) error {

	stmt := `update user set
		email = $1,
		first_name = $2,
		last_name = $3,
		user_active = $4,
		updated_at = $5
		where id = $6`

	result := db.Exec(stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Active,
		time.Now(),
		user.ID,
	)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Delete deletes one user from the database, by User.ID
func (u *User) DeleteByID(id int) error {
	result := db.Where("id = ?", id).Delete(&User{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Insert inserts a new user into the database, and returns the ID of the newly inserted row
func (u *User) Insert(user User) (uint, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}
	user.Password = string(hashedPassword)

	result := db.Create(user)
	if result.Error != nil {
		return 0, result.Error
	}

	return user.ID, nil
}

// ResetPassword is the method we will use to change a user's password.
func (u *User) ResetPassword(password string) error {
	var user *User
	var newPassword string
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}
	result := db.Model(user).Update("Password", hashedPassword)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the password
// and hash match, we return true; otherwise, we return false.
func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func (u *User) LoginUser(email string, password string) (*User, bool, string, error) {
	var user User
	query := `
			select 
			    id, 
			    email, 
			    first_name, 
			    last_name, 
			    password, 
			    user_active, 
			    is_admin, 
			    created_at, 
			    updated_at 
			from 
			    user 
			where 
			    email = $1`

	result := db.Raw(query, email).Scan(&user)
	if result.Error != nil {
		return nil, false, "", errors.New("invalid email")
	}
	//comparing and confirming the  password
	//matching the hashed password in the database to the testpassword inputed by user
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		log.Println("Incorrect Password")
		return nil, false, user.Password, errors.New("incorrect password")
	} else if err != nil {
		return nil, false, user.Password, err
	}
	var plan Plan

	/*
		get plan if any ,joins the user id and plan id into when the user selects a plan
		the user_plans table
	*/
	query = `select p.id, p.plan_name, p.plan_amount, p.created_at, p.updated_at from
			plans p
			left join userplan up on (p.id = up.plan_id)
			where up.user_id = $1`

	result = db.Raw(query, user.ID).Scan(&plan)
	if result.Error == nil {
		user.Plan = &plan
	}
	return &user, true, user.Password, nil

}
