package database

import (
  "log"
  "context"

  "backend/common"

  "github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func connectDB() {
  var err error
  DB, err = pgxpool.New(context.Background(), common.MustGetEnv("DB_URL"))
  if err != nil {
    log.Fatalf("Unable to connect to database: %v\n", err)
  }

  if err = DB.Ping(context.Background()); err != nil {
    log.Fatalf("Database connection verification failed: %v\n", err)
  }

  log.Println("Database connected successfully using pgx!")
}

