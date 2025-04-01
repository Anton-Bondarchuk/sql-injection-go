package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sql-injection-go/internal/config"
	"sql-injection-go/internal/storage/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config := config.MustLoad()

	connString := config.StorageConfig.DatabaseUrl
	storage, err := storage.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Невозможно подключиться к базе данных: %v", err)
	}

	init_table_path := "./init_table.sql"
	seeds_path := "./seeds.sql"

	mustMigrate(init_table_path, storage.Conn, "sql init table")
	mustMigrate(seeds_path, storage.Conn, "seeds")
}


func mustMigrate(path string, conn *pgxpool.Pool, op string) {
	sqlBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Ошибка чтения файла миграции: %v", err)
		panic(err)
	}

	_, err = conn.Exec(context.Background(), string(sqlBytes))
	if err != nil {
		log.Fatalf("Ошибка выполнения миграции: " + op + "%v", err)
		panic(err)
	}

	fmt.Println("Миграция " + op + " успешно выполнена!")
}	