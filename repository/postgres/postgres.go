package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/mohamadafzal06/depository/config"
	"github.com/mohamadafzal06/depository/entity"
)

var (
	ErrTableCreation = errors.New("table creation failed")
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres() (*Postgres, error) {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable",
		config.DatabaseUser, config.DatabasePass, config.DatabaseAddress, config.DatabaseDBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return &Postgres{db: db}, nil
}

func (pg *Postgres) Init() error {
	return pg.CreateAccountTable()
}

func (pg *Postgres) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
	id SERIAL PRIMARY KEY,
	firstname VARCHAR(50),
	lastname VARCHAR(50),
	encrypted_pass VARCHAR(50),
	number SERIAL UNIQUE,
	balance, 
	created_at timestamp
	CONSTRAINT number_range CHECK (number BETWEEN 10000000 AND 99999999)
	);`

	_, err := pg.db.Exec(query)
	if err != nil {
		return ErrTableCreation
	}

	return nil
}

func (pg *Postgres) CreateAccount(ctx context.Context, acc *entity.Account) (int64, error) {
	res, err := pg.db.ExecContext(ctx,
		"insert into account (firstname, lastname, encrypted_pass, balance) values($1, $2, $3, $4);",
		acc.FirstName, acc.LastName, acc.EncryptedPassword, acc.Balance)

	if err != nil {
		return -1, fmt.Errorf("cannot insert this account into db: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("error while getting inserted id: %w", err)
	}
	row := pg.db.QueryRowContext(ctx, "select number from account where id=$1", id)
	var number int64
	if err = row.Scan(&number); err != nil {
		return -1, fmt.Errorf("error while scanning number of inserted account: %w", err)
	}

	return number, nil
}

func (pg *Postgres) GetAccountByNumber(ctx context.Context, number int64) (*entity.Account, error) {
	row := pg.db.QueryRowContext(ctx, "select (fistname, lastname, balance, created_at) from account where number=$1", number)
	var acc entity.Account
	err := row.Scan(&acc.FirstName, &acc.LastName, &acc.Balance, &acc.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return &entity.Account{}, fmt.Errorf("account with this number does not exist: %w", err)
		}

		return &entity.Account{}, fmt.Errorf("error while scanning result from db: %w", err)
	}

	return &acc, nil
}

func (pg *Postgres) DeleteAccount(ctx context.Context, number int64) error {

	_, err := pg.db.ExecContext(ctx, "delete from account where number=$1", number)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("account with this number does not exist: %w", err)
		}
		return fmt.Errorf("cannot delete account by this number: %w", err)
	}

	return nil
}

func (pg *Postgres) TransferAmount(ctx context.Context, from, to, amount int64) error {
	// begin a transaction
	tx, err := pg.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// select the balance of account from
	var balance1 int64
	err = tx.QueryRowContext(ctx, "SELECT balance FROM account WHERE number = $1", from).Scan(&balance1)
	if err != nil {
		tx.Rollback()
		return err
	}

	// check that there is enough balance to transfer
	if balance1 < amount {
		tx.Rollback()
		return errors.New("insufficient balance")
	}

	// update the balance of account from
	_, err = tx.ExecContext(ctx, "UPDATE account SET balance = balance - $1 WHERE number = $2", amount, from)
	if err != nil {
		tx.Rollback()
		return err
	}

	// update the balance of account to
	_, err = tx.ExecContext(ctx, "UPDATE account SET balance = balance + $1 WHERE number = $2", amount, to)
	if err != nil {
		tx.Rollback()
		return err
	}

	// commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (pg *Postgres) AccountAuthenticity(ctx context.Context, number int64, encPass string) error {
	row := pg.db.QueryRowContext(ctx, "select encrypted_pass from account where number=$1", number)
	var trulyPass string
	err := row.Scan(&trulyPass)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("account with this number does not exist: %w", err)
		}
		return fmt.Errorf("error while scanning result from db: %w", err)
	}

	if encPass != trulyPass {
		return fmt.Errorf("the given pass is not correct: %w", err)
	}

	return nil
}
