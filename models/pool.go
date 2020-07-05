package models

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/golang/glog"
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
		glog.Fatalf("Unable to connect to database:: %v\n", err)
		os.Exit(1)
	}
	dbPool := &DbPool{pool}
	return dbPool
}
