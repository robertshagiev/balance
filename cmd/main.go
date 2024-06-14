package main

import (
	"balance/internal/config"
	"balance/internal/handler"
	"balance/internal/logger"
	"balance/internal/repository"
	"balance/internal/routers"
	"balance/internal/usecase"
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	log := logger.New()

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to load config: %v", err))
	}

	dsn := cfg.GetDSN()
	db, err := sql.Open(cfg.Database.Driver, dsn)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to connect to the database: %v", err))
	}
	defer db.Close()

	con := repository.NewRepository(db, log)
	uc := usecase.NewUsecase(con, log)
	h := handler.NewHandler(uc, log)
	r := routers.NewRouter(h)

	address := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Info(fmt.Sprintf("Server is starting at %s", address))
	if err := http.ListenAndServe(address, r); err != nil {
		log.Fatal(fmt.Sprintf("Failed to start server: %v", err))
	}
}
