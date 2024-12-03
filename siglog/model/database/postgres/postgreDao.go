package postgres

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx/v5"

	"qqweq/siglog/model/models"
)

type PostgresDao struct {
    db *pgx.Conn
    ctx context.Context
}

var (
    dao *PostgresDao
    daoOnce sync.Once
)

func GetDao(ctx context.Context) (*PostgresDao, error) {
    var err error
    if dao == nil {
	if ctx == nil {
	    ctx = context.Background()
	}
	daoOnce.Do(func() {
	    var conn *pgx.Conn
	    conn, err = connectToDb(ctx)
	    if err == nil {
		dao = &PostgresDao{ db: conn, ctx: ctx }
	    }
	})
    }
    
    return dao, err
}

func (*PostgresDao) CreateUser(user *models.User) (string, error) {
    panic("not implemented")
}

func (*PostgresDao) ReadUserByUsername(username string) (*models.User, error) {
    panic("not implemented")
}
func (*PostgresDao) DeleteUser(user *models.User) error {
    panic("not implemented")
} 

func (*PostgresDao) CreateSession(username string) (string, error) {
    panic("not implemented")
}
func (*PostgresDao) DeleteSession(sessionId string) error {
    panic("not implemented")
}
func (*PostgresDao) UsernameFromSessionId(sessionId string) (string, error) {
    panic("not implemented")
}


func connectToDb(ctx context.Context) (*pgx.Conn, error)  { 
    connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

    conn, err := pgx.Connect(ctx, connString)
    if err != nil {
	err = fmt.Errorf("Failed to connect to the db. %w", err)
	return nil, err 
    }

    if err := conn.Ping(ctx); err != nil {
	err = fmt.Errorf("Failed to connect to the db. %w", err)
	return nil, err 
    }

    return conn, nil
}
