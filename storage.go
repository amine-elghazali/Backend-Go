package main

import (
	"database/sql"
	"fmt"

	models "github.com/amine-elghazali/Backend-Go/models"

	_ "github.com/lib/pq"
)

type Account = models.Account

type Storage interface {
	GetAccounts() ([]*Account, error)
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=backend sslmode=disable"

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) init() error {
	return s.CreateAccountTable()
}

func (s *PostgresStore) CreateAccountTable() error {

	query := `CREATE TABLE IF NOT EXISTS account (
		id SERIAL PRIMARY KEY,
		first_name varchar(50),
		last_name varchar(50),
		number serial,
		balance serial,
		created_at timestamp )
	`
	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) CreateAccount(account *Account) error {

	sqlStatement := `insert into account 
	(first_name,last_name,number,balance,created_at)
	values($1,$2,$3,$4,$5)
	`
	resp, err := s.db.Exec(
		sqlStatement,
		account.FirstName,
		account.LastName,
		account.Number,
		account.Balance,
		account.CreatedAt,
	)

	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)

	return nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {

	_, err := s.db.Query("DELETE FROM account where id=$1", id)
	return err
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {

	sqlStatement := `SELECT * from account where id=$1`

	row, err := s.db.Query(sqlStatement, id)

	if err != nil {
		return nil, err
	}

	for row.Next() {
		return scanIntoAnAccount(row)
	}

	return nil, fmt.Errorf("Account %d not found", id)

}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {

	rows, err := s.db.Query("SELECT * from account")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	accounts := []*Account{}

	for rows.Next() {

		account, err := scanIntoAnAccount(rows)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)

	}

	return accounts, nil
}

func scanIntoAnAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)

	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return account, nil
}
