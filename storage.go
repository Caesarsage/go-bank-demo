package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountById(int) (*Account, error)
	GetAccounts() ([]*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {

	db_url := os.Getenv("DATABASE_URL")

	fmt.Println(db_url)

	connStr := "user=postgres dbname=go-bank password=admin sslmode=disable"

	db, err := sql.Open("postgres", connStr)

	if err !=nil {
		return nil, err
	}

	if err := db.Ping(); err !=nil {
		return nil, err
	}

	return &PostgresStore{
		db:db,
	}, nil
}

func (s *PostgresStore) init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
				id serial primary key,
				first_name varchar(50),
				last_name varchar(50),
				number serial,
				balance serial,
				created_at timestamp
			)`

	_, err := s.db.Exec(query)
	return err
}


func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `
		INSERT INTO account (first_name, last_name, number, balance, created_at)
	 	VALUES ($1, $2, $3, $4, $5)
	`
	resp, err := s.db.Query(
		query, 
		acc.FirstName,
		acc.LastName, 
		acc.Number, 
		acc.Balance, 
		acc.CreatedAt)

	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)
	
	return nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	query := `SELECT * FROM account`

	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for rows.Next() {
		acc, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, acc)

	}

	return accounts, nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Query(`DELETE FROM account WHERE id = $1`, id)

	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	rows, err := s.db.Query(`SELECT * FROM account WHERE id = $1`, id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account not found")
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
		acc := &Account{}
		err := rows.Scan(
			&acc.ID, 
			&acc.FirstName, 
			&acc.LastName, 
			&acc.Number, 
			&acc.Balance, 
			&acc.CreatedAt,
		);

		return acc, err

}
