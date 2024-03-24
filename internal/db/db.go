package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"sso_service_grps/internal/domain/models"
	"time"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(user, password, addr, dbName string, port int) (*Storage, error) {
	var storage Storage

	const defaultMaxConns = int32(4)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	conf, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, addr, port, dbName))
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	conf.MaxConns = defaultMaxConns
	conf.MinConns = defaultMinConns
	conf.MaxConnLifetime = defaultMaxConnLifetime
	conf.MaxConnIdleTime = defaultMaxConnIdleTime
	conf.HealthCheckPeriod = defaultHealthCheckPeriod
	conf.ConnConfig.ConnectTimeout = defaultConnectTimeout

	conf.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		log.Println("Before acquiring the connection pool to the database!!")
		return true
	}

	conf.AfterRelease = func(c *pgx.Conn) bool {
		log.Println("After releasing the connection pool to the database!!")
		return true
	}

	conf.BeforeClose = func(c *pgx.Conn) {
		log.Println("Closed the connection pool to the database!!")
	}

	connPool, err := pgxpool.NewWithConfig(context.Background(), conf)
	if err != nil {
		log.Fatal("Error while creating connection to the database!!")
	}

	connection, err := connPool.Acquire(context.Background())
	if err != nil {
		log.Fatal("Error while acquiring connection from the database pool!!")
	}
	defer connection.Release()
	err = connection.Ping(context.Background())
	if err != nil {
		log.Fatal("Could not ping database")
	}
	fmt.Println("Connected to the database!!")

	storage.db = connPool

	return &storage, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, password string) (int64, error) {
	query := "INSERT INTO users(email, password) VALUES ($1,$2) returning id"
	res, err := s.db.Query(context.Background(), query, email, password)
	if err != nil {
		return 0, err
	}
	ind, err := pgx.CollectOneRow(res, pgx.RowTo[int])
	if err != nil {
		return 0, err
	}
	res.Close()
	return int64(ind), nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	query := "SELECT * FROM users WHERE email = $1"
	rows, err := s.db.Query(context.Background(), query, email)
	if err != nil {
		return models.User{}, err
	}
	defer rows.Close()
	res, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return models.User{}, err
	}
	return res, nil
}

func (s *Storage) App(ctx context.Context, id int) (models.App, error) {
	query := "SELECT * FROM apps WHERE id = $1"
	rows, err := s.db.Query(context.Background(), query, id)
	if err != nil {
		return models.App{}, err
	}
	defer rows.Close()
	res, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.App])
	if err != nil {
		return models.App{}, err
	}
	return res, nil
}
