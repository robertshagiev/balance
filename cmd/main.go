package main

import (
	"balance/internal/handler"
	"balance/internal/repository"
	"balance/internal/routers"
	"balance/internal/usecase"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "root:root@tcp(localhost:3306)/balance?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	con := repository.NewRepository(db)
	uc := usecase.NewUsecase(con)
	h := handler.NewHandler(uc)

	r := routers.NewRouter(h)

	log.Println("Server is starting at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
