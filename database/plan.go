package database

import (
	"fmt"

	"gorm.io/gorm"
)

// Plan is the type for subscription plans
type Plan struct {
	gorm.Model
	PlanName            string
	PlanAmount          int
	PlanAmountFormatted string    `gorm:"-"`
	UserPlan            *UserPlan `gorm:"foreignKey:UserID"`
}

func (p *Plan) GetAll() ([]*Plan, error) {
	var plans []*Plan

	result := db.Order("id").Find(&plans)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, plan := range plans {
		plan.PlanAmountFormatted = plan.AmountForDisplay()
	}
	return plans, nil
}

// GetOne returns one plan by id
func (p *Plan) GetOne(id int) (*Plan, error) {
	//when user slects a plan the id in the plan db is saved into
	//the user_plan db
	var plan Plan
	result := db.First(&plan, id)
	if result.Error != nil {
		return nil, result.Error
	}
	plan.PlanAmountFormatted = plan.AmountForDisplay()

	return &plan, nil
}

// SubscribeUserToPlan subscribes a user to one plan by insert
// values into userplan table
func (p *Plan) SubscribeUserToPlan(user User, plan Plan) error {
	tx := db.Begin() // Start a transaction

	// Delete existing plan for the user
	if err := tx.Where("user_id = ?", user.ID).Delete(&UserPlan{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Subscribe user to the new plan
	newUserPlan := &UserPlan{
		UserID: user.ID,
		PlanID: plan.ID,
	}

	if err := tx.Create(newUserPlan).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// AmountForDisplay formats the price we have in the DB as a currency string
func (p *Plan) AmountForDisplay() string {
	amount := float64(p.PlanAmount) / 100.0
	return fmt.Sprintf("$%.2f", amount)
}
