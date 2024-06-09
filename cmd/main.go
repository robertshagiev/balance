package main

import (
	"balance/internal/handler"
	"balance/internal/repository"
	"balance/internal/server"
	"balance/internal/usecase"
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "root:root@tcp(localhost:3306)/balance_service?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)
	uc := usecase.NewUsecase(repo)
	h := handler.NewHandler(uc)

	srv := server.NewServer(h, "localhost", "8080")
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
