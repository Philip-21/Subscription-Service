package database

import (
	"context"
	"database/sql"
	"log"
	"time"
)

func Users(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS subscription_users(
		id Serial primary key,
		email character varying(255) NOT NULL,
		first_name character varying(255) NOT NULL,
		last_name character varying(255) NOT NULL,
		password character varying(60) NOT NULL,
		user_active int DEFAULT 0,
		is_admin int DEFAULT 0,
		created_at timestamp(6) NOT NULL, 
		updated_at timestamp(6) NOT NULL
	)`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s in Creating subscription_users Table", err)
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
		return err
	}
	log.Println("subscription_users Table Created")
	return nil
}

func User_Plans(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS user_plans(
		id Serial primary key,
		user_id int NOT NULL,
		plan_id int NOT NULL,
		created_at timestamp(6) NOT NULL,
		updated_at timestamp(6) NOT NULL

	)`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	uplan, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s in creating User-Plans table", err)
		return err
	}
	_, err = uplan.RowsAffected()
	if err != nil {
		log.Printf("Error %s when geting rows affected", err)
		return err
	}
	log.Println("User-Plans Table Created")
	return nil
}

func Plans(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS plans(
		id Serial primary key,
		plan_name character varying(255) NOT NULL,
		plan_amount int NOT NULL,
		created_at timestamp(6) NOT NULL, 
		updated_at timestamp(6) NOT NULL
		
	)`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	plan, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s in creating Plans table", err)
		return err
	}
	_, err = plan.RowsAffected()
	if err != nil {
		log.Printf("Error %s when geting rows affected", err)
		return err
	}
	log.Println("Plans Table Created")
	return nil

}

func Altertable(db *sql.DB) error {
	query := `ALTER TABLE user_plans ADD FOREIGN KEY (user_id) REFERENCES subscription_users(id);
	          ALTER TABLE user_plans ADD FOREIGN KEY (plan_id) REFERENCES plans(id)`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	alter, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s in altering User-Plans table", err)
		return err
	}
	_, err = alter.RowsAffected()
	if err != nil {
		log.Printf("Error %s when geting rows affected", err)
		return err
	}
	log.Println("User-Plan Table Altered")
	return nil

}

func SeedUserPlans(db *sql.DB, p *Plan) error {
	query := `INSERT INTO plans (id, plan_name, plan_amount, created_at, updated_at)
	      VALUES($1, $2, $3, $4, $5)
		ON CONFLICT(id) DO NOTHING;`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, p.ID, p.PlanName, p.PlanAmount, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		log.Printf("Error %s when inserting row into rooms table", err)
		log.Println("seeding db affected")
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Println("plan value inserted")
	return nil

}
