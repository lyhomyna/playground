package postgres

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx/v5"

	"qqweq/siglog/model/database"
	"qqweq/siglog/model/models"
)

type PostgresDao struct {
    db *pgx.Conn
}

var (
    dao *PostgresDao
    daoOnce sync.Once
)

func NewDao(ctx context.Context) (*PostgresDao, error) {
    if ctx == nil {
	ctx = context.Background()
    }

    var err error
    daoOnce.Do(func() {
	var conn *pgx.Conn
	conn, err = connectToDb(ctx)
	if err == nil {
	    dao = &PostgresDao{ conn }
	}
    })
    
    return dao, err
}

func (*PostgresDao) CreateUser(user *models.User) (string, error) {
    // STUB
    return "STUB", nil
}

func (*PostgresDao) ReadUserByUsername(username string) (*models.User, error) {
    // STUB
    return nil, nil
}
func (*PostgresDao) DeleteUser(user *models.User) error {
    // STUB
    return nil 
} 

func (*PostgresDao) CreateSession(username string) (string, error) {
    // STUB
    return "STUB", nil 
}
func (*PostgresDao) DeleteSession(sessionId string) error {
    // STUB
    return nil 
}
func (*PostgresDao) UsernameFromSessionId(sessionId string) (string, error) {
    // STUB
    return "STUB", nil
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

// Validate if the SiglogDao interface is implemented by PostgresDao
var _ database.SiglogDao = (*PostgresDao)(nil)
