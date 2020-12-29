package models

import (
	"context"
	"os"

	"github.com/golang/glog"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DbPool struct {
	db *pgxpool.Pool
}

func CreatePool(dbUrl string) *DbPool {
	poolConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		glog.Fatalf("Error in DB Configuration:: %v\n", err)
		os.Exit(1)
	}
	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		glog.Fatalf("Unable to connect to database (make sure db is up and env.sh is correct):: %v\n", err)
		os.Exit(1)
	}
	dbPool := &DbPool{pool}
	return dbPool
}
